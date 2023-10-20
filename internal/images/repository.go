package images

type Repository interface {
	Create(imgName string, imgData []byte) (location string, err error)
	Get(location string) (imgData []byte, err error)
	Update(location string, imgData []byte) (err error)
	Delete(location string) (err error)
}
