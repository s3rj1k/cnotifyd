package main

// ContainerDB contains operational container data.
type ContainerDB struct {
	mountBidiName *BidiIntToStrDict
	pidBidiName   *BidiIntToStrDict
}

// NewContainerDB creates new NewContainerDB object.
func NewContainerDB() *ContainerDB {
	db := new(ContainerDB)
	db.mountBidiName = NewBidiIntToStrDict()
	db.pidBidiName = NewBidiIntToStrDict()

	return db
}

// BindMountIDToContainerName binds container rootfs mount ID to a container name for bidirectional lookup.
func (db *ContainerDB) BindMountIDToContainerName(mountID int, name string) {
	db.mountBidiName.Bind(mountID, name)
}

// BindPIDToContainerName binds container PID to a container name for bidirectional lookup.
func (db *ContainerDB) BindPIDToContainerName(pid int, name string) {
	db.pidBidiName.Bind(pid, name)
}

// UnBind is used to undefine specified key.
func (db *ContainerDB) UnBind(keys ...interface{}) {
	db.mountBidiName.UnBind(keys)
	db.pidBidiName.UnBind(keys)
}

// GetMountIDFromName returns rootfs mount ID for specified container name with bool that specifies successful lookup.
func (db *ContainerDB) GetMountIDFromName(name string) (int, bool) {
	return db.mountBidiName.GetInteger(name)
}

// GetContainerNameFromMountID returns container name for specified mount ID with bool that specifies successful lookup.
func (db *ContainerDB) GetContainerNameFromMountID(id int) (string, bool) {
	return db.mountBidiName.GetString(id)
}

// GetPIDFromName returns Init PID of specified container name with bool that specifies successful lookup.
func (db *ContainerDB) GetPIDFromName(name string) (int, bool) {
	return db.pidBidiName.GetInteger(name)
}

// GetNameFromPID returns container name for specified Init PID with bool that specifies successful lookup.
func (db *ContainerDB) GetNameFromPID(pid int) (string, bool) {
	return db.pidBidiName.GetString(pid)
}

// GetMountIDFromPID returns rootfs mount ID for specified container Init PID with bool that specifies successful lookup.
func (db *ContainerDB) GetMountIDFromPID(pid int) (int, bool) {
	if name, ok := db.pidBidiName.GetString(pid); ok {
		return db.mountBidiName.GetInteger(name)
	}

	return 0, false
}

// GetPIDFromMountID returns container PID from specified container rootfs mount ID with bool that specifies successful lookup.
func (db *ContainerDB) GetPIDFromMountID(id int) (int, bool) {
	if name, ok := db.mountBidiName.GetString(id); ok {
		return db.pidBidiName.GetInteger(name)
	}

	return 0, false
}
