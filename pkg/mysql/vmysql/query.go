/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreedto in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vmysql

import (
	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/proto/vtrpc"
	"vitess.io/vitess/go/vt/vterrors"

	querypb "vitess.io/vitess/go/vt/proto/query"
)

// This file contains the methods related to queries.

//
// Client side methods.
//

// WriteComQuery writes a query for the server to execute.
// Client -> Server.
// Returns SQLError(CRServerGone) if it can't.
func (c *Conn) WriteComQuery(query string) error {
	// This is a new command, need to reset the sequence.
	c.sequence = 0

	data := c.startEphemeralPacket(len(query) + 1)
	data[0] = ComQuery
	copy(data[1:], query)
	if err := c.writeEphemeralPacket(); err != nil {
		return NewSQLError(CRServerGone, SSUnknownSQLState, err.Error())
	}
	return nil
}

// readColumnDefinition reads the next Column Definition packet.
// Returns a SQLError.
func (c *Conn) readColumnDefinition(field *querypb.Field, index int) error {
	colDef, err := c.readEphemeralPacket()
	if err != nil {
		return NewSQLError(CRServerLost, SSUnknownSQLState, "%v", err)
	}
	defer c.RecycleReadPacket()

	// Catalog is ignored, always set to "def"
	pos, ok := skipLenEncString(colDef, 0)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "skipping col %v catalog failed", index)
	}

	// schema, table, orgTable, name and OrgName are strings.
	field.Database, pos, ok = readLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v schema failed", index)
	}
	field.Table, pos, ok = readLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v table failed", index)
	}
	field.OrgTable, pos, ok = readLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v org_table failed", index)
	}
	field.Name, pos, ok = readLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v name failed", index)
	}
	field.OrgName, pos, ok = readLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v org_name failed", index)
	}

	// Skip length of fixed-length fields.
	pos++

	// characterSet is a uint16.
	characterSet, pos, ok := readUint16(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v characterSet failed", index)
	}
	field.Charset = uint32(characterSet)

	// columnLength is a uint32.
	field.ColumnLength, pos, ok = readUint32(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v columnLength failed", index)
	}

	// type is one byte.
	t, pos, ok := readByte(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v type failed", index)
	}

	// flags is 2 bytes.
	flags, pos, ok := readUint16(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v flags failed", index)
	}

	// Convert MySQL type to Vitess type.
	field.Type, err = sqltypes.MySQLToType(int64(t), int64(flags))
	if err != nil {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "MySQLToType(%v,%v) failed for column %v: %v", t, flags, index, err)
	}

	// Decimals is a byte.
	decimals, _, ok := readByte(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v decimals failed", index)
	}
	field.Decimals = uint32(decimals)

	// If we didn't get column length or character set,
	// we assume the orignal row on the other side was encoded from
	// a Field without that data, so we don't return the flags.
	if field.ColumnLength != 0 || field.Charset != 0 {
		field.Flags = uint32(flags)

		// FIXME(alainjobart): This is something the MySQL
		// client library does: If the type is numerical, it
		// adds a NUM_FLAG to the flags.  We're doing it here
		// only to be compatible with the C library. Once
		// we're not using that library any more, we'll remove this.
		// See doc.go.
		if IsNum(t) {
			field.Flags |= uint32(querypb.MySqlFlag_NUM_FLAG)
		}
	}

	return nil
}

// readColumnDefinitionType is a faster version of
// readColumnDefinition that only fills in the Type.
// Returns a SQLError.
func (c *Conn) readColumnDefinitionType(field *querypb.Field, index int) error {
	colDef, err := c.readEphemeralPacket()
	if err != nil {
		return NewSQLError(CRServerLost, SSUnknownSQLState, "%v", err)
	}
	defer c.RecycleReadPacket()

	// catalog, schema, table, orgTable, name and orgName are
	// strings, all skipped.
	pos, ok := skipLenEncString(colDef, 0)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "skipping col %v catalog failed", index)
	}
	pos, ok = skipLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "skipping col %v schema failed", index)
	}
	pos, ok = skipLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "skipping col %v table failed", index)
	}
	pos, ok = skipLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "skipping col %v org_table failed", index)
	}
	pos, ok = skipLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "skipping col %v name failed", index)
	}
	pos, ok = skipLenEncString(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "skipping col %v org_name failed", index)
	}

	// Skip length of fixed-length fields.
	pos++

	// characterSet is a uint16.
	_, pos, ok = readUint16(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v characterSet failed", index)
	}

	// columnLength is a uint32.
	_, pos, ok = readUint32(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v columnLength failed", index)
	}

	// type is one byte
	t, pos, ok := readByte(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v type failed", index)
	}

	// flags is 2 bytes
	flags, _, ok := readUint16(colDef, pos)
	if !ok {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extracting col %v flags failed", index)
	}

	// Convert MySQL type to Vitess type.
	field.Type, err = sqltypes.MySQLToType(int64(t), int64(flags))
	if err != nil {
		return NewSQLError(CRMalformedPacket, SSUnknownSQLState, "MySQLToType(%v,%v) failed for column %v: %v", t, flags, index, err)
	}

	// skip decimals

	return nil
}

// parseRow parses an individual row.
// Returns a SQLError.
func (c *Conn) parseRow(data []byte, fields []*querypb.Field) ([]sqltypes.Value, error) {
	colNumber := len(fields)
	result := make([]sqltypes.Value, colNumber)
	pos := 0
	for i := 0; i < colNumber; i++ {
		if data[pos] == 0xfb {
			pos++
			continue
		}
		var s []byte
		var ok bool
		s, pos, ok = readLenEncStringAsBytes(data, pos)
		if !ok {
			return nil, NewSQLError(CRMalformedPacket, SSUnknownSQLState, "decoding string failed")
		}
		result[i] = sqltypes.MakeTrusted(fields[i].Type, s)
	}
	return result, nil
}

// ExecuteFetch executes a query and returns the result.
// Returns a SQLError. Depending on the transport used, the error
// returned might be different for the same condition:
//
// 1. if the server closes the connection when no command is in flight:
//
//		1.1 unix: WriteComQuery will fail with a 'broken pipe', and we'll
//		    return CRServerGone(2006).
//
//		1.2 tcp: WriteComQuery will most likely work, but readComQueryResponse
//		    will fail, and we'll return CRServerLost(2013).
//
//		    This is because closing a TCP socket on the server side sends
//		    a FIN to the client (telling the client the server is done
//		    writing), but on most platforms doesn't send a RST.  So the
//		    client has no idea it can't write. So it succeeds writing data, which
//		    *then* triggers the server to send a RST back, received a bit
//		    later. By then, the client has already started waiting for
//		    the response, and will just return a CRServerLost(2013).
//		    So CRServerGone(2006) will almost never be seen with TCP.
//
//	 2. if the server closes the connection when a command is in flight,
//	    readComQueryResponse will fail, and we'll return CRServerLost(2013).
func (c *Conn) ExecuteFetch(query string, maxrows int, wantfields bool) (result *sqltypes.Result, err error) {
	result, _, err = c.ExecuteFetchMulti(query, maxrows, wantfields)
	return result, err
}

// ExecuteFetchMulti is for fetching multiple results from a multi-statement result.
// It returns an additional 'more' flag. If it is set, you must fetch the additional
// results using ReadQueryResult.
func (c *Conn) ExecuteFetchMulti(query string, maxrows int, wantfields bool) (result *sqltypes.Result, more bool, err error) {
	defer func() {
		if err != nil {
			if sqlerr, ok := err.(*SQLError); ok {
				sqlerr.Query = query
			}
		}
	}()

	// Send the query as a COM_QUERY packet.
	if err = c.WriteComQuery(query); err != nil {
		return nil, false, err
	}

	res, more, _, err := c.ReadQueryResult(maxrows, wantfields)
	return res, more, err
}

// ExecuteFetchWithWarningCount is for fetching results and a warning count
// Note: In a future iteration this should be abolished and merged into the
// ExecuteFetch API.
func (c *Conn) ExecuteFetchWithWarningCount(query string, maxrows int, wantfields bool) (result *sqltypes.Result, warnings uint16, err error) {
	defer func() {
		if err != nil {
			if sqlerr, ok := err.(*SQLError); ok {
				sqlerr.Query = query
			}
		}
	}()

	// Send the query as a COM_QUERY packet.
	if err = c.WriteComQuery(query); err != nil {
		return nil, 0, err
	}

	res, _, warnings, err := c.ReadQueryResult(maxrows, wantfields)
	return res, warnings, err
}

// ReadQueryResult gets the result from the last written query.
func (c *Conn) ReadQueryResult(maxrows int, wantfields bool) (result *sqltypes.Result, more bool, warnings uint16, err error) {
	// Get the result.
	affectedRows, lastInsertID, colNumber, more, warnings, err := c.readComQueryResponse()
	if err != nil {
		return nil, false, 0, err
	}

	if colNumber == 0 {
		// OK packet, means no results. Just use the numbers.
		return &sqltypes.Result{
			RowsAffected: affectedRows,
			InsertID:     lastInsertID,
		}, more, warnings, nil
	}

	fields := make([]querypb.Field, colNumber)
	result = &sqltypes.Result{
		Fields: make([]*querypb.Field, colNumber),
	}

	// Read column headers. One packet per column.
	// Build the fields.
	for i := 0; i < colNumber; i++ {
		result.Fields[i] = &fields[i]

		if wantfields {
			if err := c.readColumnDefinition(result.Fields[i], i); err != nil {
				return nil, false, 0, err
			}
		} else {
			if err := c.readColumnDefinitionType(result.Fields[i], i); err != nil {
				return nil, false, 0, err
			}
		}
	}

	if c.Capabilities&CapabilityClientDeprecateEOF == 0 {
		// EOF is only present here if it's not deprecated.
		data, err := c.readEphemeralPacket()
		if err != nil {
			return nil, false, 0, NewSQLError(CRServerLost, SSUnknownSQLState, "%v", err)
		}
		if isEOFPacket(data) {

			// This is what we expect.
			// Warnings and status flags are ignored.
			c.RecycleReadPacket()
			// goto: read row loop

		} else if isErrorPacket(data) {
			defer c.RecycleReadPacket()
			return nil, false, 0, ParseErrorPacket(data)
		} else {
			defer c.RecycleReadPacket()
			return nil, false, 0, vterrors.Errorf(vtrpc.Code_INTERNAL, "unexpected packet after fields: %v", data)
		}
	}

	// read each row until EOF or OK packet.
	for {
		data, err := c.ReadPacket()
		if err != nil {
			return nil, false, 0, err
		}

		if isEOFPacket(data) {
			// Strip the partial Fields before returning.
			if !wantfields {
				result.Fields = nil
			}
			result.RowsAffected = uint64(len(result.Rows))

			// The deprecated EOF packets change means that this is either an
			// EOF packet or an OK packet with the EOF type code.
			if c.Capabilities&CapabilityClientDeprecateEOF == 0 {
				warnings, more, err = parseEOFPacket(data)
				if err != nil {
					return nil, false, 0, err
				}
			} else {
				var statusFlags uint16
				_, _, statusFlags, warnings, err = parseOKPacket(data)
				if err != nil {
					return nil, false, 0, err
				}
				more = (statusFlags & ServerMoreResultsExists) != 0
			}
			return result, more, warnings, nil

		} else if isErrorPacket(data) {
			// Error packet.
			return nil, false, 0, ParseErrorPacket(data)
		}

		// Check we're not over the limit before we add more.
		if len(result.Rows) == maxrows {
			if err := c.drainResults(); err != nil {
				return nil, false, 0, err
			}
			return nil, false, 0, NewSQLError(ERVitessMaxRowsExceeded, SSUnknownSQLState, "Row count exceeded %d", maxrows)
		}

		// Regular row.
		row, err := c.parseRow(data, result.Fields)
		if err != nil {
			return nil, false, 0, err
		}
		result.Rows = append(result.Rows, row)
	}
}

// drainResults will read all packets for a result set and ignore them.
func (c *Conn) drainResults() error {
	for {
		data, err := c.readEphemeralPacket()
		if err != nil {
			return NewSQLError(CRServerLost, SSUnknownSQLState, "%v", err)
		}
		if isEOFPacket(data) {
			c.RecycleReadPacket()
			return nil
		} else if isErrorPacket(data) {
			defer c.RecycleReadPacket()
			return ParseErrorPacket(data)
		}
		c.RecycleReadPacket()
	}
}

func (c *Conn) readComQueryResponse() (affectedRows uint64, lastInsertID uint64, status int, more bool, warnings uint16, err error) {
	data, err := c.readEphemeralPacket()
	if err != nil {
		return 0, 0, 0, false, 0, NewSQLError(CRServerLost, SSUnknownSQLState, "%v", err)
	}
	defer c.RecycleReadPacket()
	if len(data) == 0 {
		return 0, 0, 0, false, 0, NewSQLError(CRMalformedPacket, SSUnknownSQLState, "invalid empty COM_QUERY response packet")
	}

	switch data[0] {
	case OKPacket:
		affectedRows, lastInsertID, status, warnings, err := parseOKPacket(data)
		return affectedRows, lastInsertID, 0, (status & ServerMoreResultsExists) != 0, warnings, err
	case ErrPacket:
		// Error
		return 0, 0, 0, false, 0, ParseErrorPacket(data)
	case 0xfb:
		// Local infile
		return 0, 0, 0, false, 0, vterrors.Errorf(vtrpc.Code_UNIMPLEMENTED, "not implemented")
	}
	n, pos, ok := readLenEncInt(data, 0)
	if !ok {
		return 0, 0, 0, false, 0, NewSQLError(CRMalformedPacket, SSUnknownSQLState, "cannot get column number")
	}
	if pos != len(data) {
		return 0, 0, 0, false, 0, NewSQLError(CRMalformedPacket, SSUnknownSQLState, "extra data in COM_QUERY response")
	}
	return 0, 0, int(n), false, 0, nil
}

//
// Server side methods.
//

func (c *Conn) parseComQuery(data []byte) string {
	return string(data[1:])
}

func (c *Conn) parseComSetOption(data []byte) (uint16, bool) {
	val, _, ok := readUint16(data, 1)
	return val, ok
}

func (c *Conn) parseComInitDB(data []byte) string {
	return string(data[1:])
}

func (c *Conn) sendColumnCount(count uint64) error {
	length := lenEncIntSize(count)
	data := c.startEphemeralPacket(length)
	writeLenEncInt(data, 0, count)
	return c.writeEphemeralPacket()
}

func (c *Conn) writeColumnDefinition(field *querypb.Field) error {
	length := 4 + // lenEncStringSize("def")
		lenEncStringSize(field.Database) +
		lenEncStringSize(field.Table) +
		lenEncStringSize(field.OrgTable) +
		lenEncStringSize(field.Name) +
		lenEncStringSize(field.OrgName) +
		1 + // length of fixed length fields
		2 + // character set
		4 + // column length
		1 + // type
		2 + // flags
		1 + // decimals
		2 // filler

	// Get the type and the flags back. If the Field contains
	// non-zero flags, we use them. Otherwise use the flags we
	// derive from the type.
	typ, flags := sqltypes.TypeToMySQL(field.Type)
	if field.Flags != 0 {
		flags = int64(field.Flags)
	}

	data := c.startEphemeralPacket(length)
	pos := 0

	pos = writeLenEncString(data, pos, "def") // Always the same.
	pos = writeLenEncString(data, pos, field.Database)
	pos = writeLenEncString(data, pos, field.Table)
	pos = writeLenEncString(data, pos, field.OrgTable)
	pos = writeLenEncString(data, pos, field.Name)
	pos = writeLenEncString(data, pos, field.OrgName)
	pos = writeByte(data, pos, 0x0c)
	pos = writeUint16(data, pos, uint16(field.Charset))
	pos = writeUint32(data, pos, field.ColumnLength)
	pos = writeByte(data, pos, byte(typ))
	pos = writeUint16(data, pos, uint16(flags))
	pos = writeByte(data, pos, byte(field.Decimals))
	pos = writeUint16(data, pos, uint16(0x0000))

	if pos != len(data) {
		return vterrors.Errorf(vtrpc.Code_INTERNAL, "packing of column definition used %v bytes instead of %v", pos, len(data))
	}

	return c.writeEphemeralPacket()
}

func (c *Conn) writeRow(row []sqltypes.Value) error {
	length := 0
	for _, val := range row {
		if val.IsNull() {
			length++
		} else {
			l := len(val.Raw())
			length += lenEncIntSize(uint64(l)) + l
		}
	}

	data := c.startEphemeralPacket(length)
	pos := 0
	for _, val := range row {
		if val.IsNull() {
			pos = writeByte(data, pos, NullValue)
		} else {
			l := len(val.Raw())
			pos = writeLenEncInt(data, pos, uint64(l))
			pos += copy(data[pos:], val.Raw())
		}
	}

	if pos != length {
		return vterrors.Errorf(vtrpc.Code_INTERNAL, "packet row: got %v bytes but expected %v", pos, length)
	}

	return c.writeEphemeralPacket()
}

// writeFields writes the fields of a Result. It should be called only
// if there are valid columns in the result.
func (c *Conn) writeFields(result *sqltypes.Result) error {
	// Send the number of fields first.
	if err := c.sendColumnCount(uint64(len(result.Fields))); err != nil {
		return err
	}

	// Now send each Field.
	for _, field := range result.Fields {
		if err := c.writeColumnDefinition(field); err != nil {
			return err
		}
	}

	// Now send an EOF packet.
	if c.Capabilities&CapabilityClientDeprecateEOF == 0 {
		// With CapabilityClientDeprecateEOF, we do not send this EOF.
		if err := c.writeEOFPacket(c.StatusFlags, 0); err != nil {
			return err
		}
	}
	return nil
}

// writeRows sends the rows of a Result.
func (c *Conn) writeRows(result *sqltypes.Result) error {
	for _, row := range result.Rows {
		if err := c.writeRow(row); err != nil {
			return err
		}
	}
	return nil
}

// writeEndResult concludes the sending of a Result.
// if more is set to true, then it means there are more results afterwords
func (c *Conn) writeEndResult(more bool, affectedRows, lastInsertID uint64, warnings uint16) error {
	// Send either an EOF, or an OK packet.
	// See doc.go.
	flags := c.StatusFlags
	if more {
		flags |= ServerMoreResultsExists
	}
	if c.Capabilities&CapabilityClientDeprecateEOF == 0 {
		if err := c.writeEOFPacket(flags, warnings); err != nil {
			return err
		}
	} else {
		// This will flush too.
		if err := c.writeOKPacketWithEOFHeader(affectedRows, lastInsertID, flags, warnings); err != nil {
			return err
		}
	}

	return nil
}
