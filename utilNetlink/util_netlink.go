package utilNetlink

type UtilNetlink struct {
}

func NewUtilNetlink() *UtilNetlink {
	return &UtilNetlink{}
}
func (u *UtilNetlink) SrcRouteAdd(src string, gateway string) (err error) {
	return SrcRouteAdd(src, gateway)
}
func (u *UtilNetlink) SrcRouteDelete(src string) (err error) {
	return SrcRouteDelete(src)
}
