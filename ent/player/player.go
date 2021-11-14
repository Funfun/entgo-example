// Code generated by entc, DO NOT EDIT.

package player

const (
	// Label holds the string label denoting the player type in the database.
	Label = "player"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldNickname holds the string denoting the nickname field in the database.
	FieldNickname = "nickname"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldScores holds the string denoting the scores field in the database.
	FieldScores = "scores"
	// Table holds the table name of the player in the database.
	Table = "players"
)

// Columns holds all SQL columns for player fields.
var Columns = []string{
	FieldID,
	FieldNickname,
	FieldEmail,
	FieldScores,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}