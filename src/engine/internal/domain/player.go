package domain

// player base data
type Player struct {
	ID    string
	Name  string
	Role  Role
	Alive bool
	// TODO: add personallity trait (e.g. timid, agressive, nuetral...)
}

// possible player roles
type Role int

const (
	RoleUnknown Role = iota
	RoleVillager
	RoleMafia
	RoleDoctor
	RoleSheriff
)

func (r Role) String() string {
	switch r {
	case RoleUnknown:
		return "unknown"
	case RoleVillager:
		return "villager"
	case RoleMafia:
		return "mafia"
	case RoleDoctor:
		return "doctor"
	case RoleSheriff:
		return "sheriff"
	default:
		return "invalid"
	}
}

// player state helpers
func (p Player) IsAlive() bool {
	return p.Alive
}

// player role helpers
func (r Role) IsVillagerTeam() bool {
	return r == RoleVillager ||
		r == RoleDoctor ||
		r == RoleSheriff
}

func (r Role) IsMafiaTeam() bool {
	return r == RoleMafia
}

func (r Role) HasNightAction() bool {
	return r == RoleMafia ||
		r == RoleDoctor ||
		r == RoleSheriff
}
