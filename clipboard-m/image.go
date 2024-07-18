package clipboardm

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"runtime"
	"unsafe"
)

type BITMAPINFOHEADER struct {
	biSize          uint32
	biWidth         uint32
	biHeight        uint32
	biPlanes        uint16
	biBitCount      uint16
	biCompression   uint32
	biSizeImage     uint32
	biXPelsPerMeter int32
	biYPelsPerMeter int32
	biClrUsed       uint32
	biClrImportant  uint32
}

func ReadImage() ([]byte, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if formatAvailable, _, err := isClipboardFormatAvailable.Call(CF_DIB); formatAvailable == 0 {
		return []byte{}, err
	}
	err := waitOpenClipboard()
	if err != nil {
		return []byte{}, err
	}

	h, _, err := getClipboardData.Call(CF_DIB)
	if h == 0 {
		_, _, _ = closeClipboard.Call()
		return []byte{}, err
	}

	l, _, err := globalLock.Call(h)
	if l == 0 {
		_, _, _ = closeClipboard.Call()
		return []byte{}, err
	}

	// ls, _, err := globalSize.Call(l)

	dibHeader := (*BITMAPINFOHEADER)(unsafe.Pointer(l))
	if dibHeader.biSize != uint32(unsafe.Sizeof(BITMAPINFOHEADER{})) {
		return nil, fmt.Errorf("invalid DIB header size")
	}

	offset := int(uintptr(dibHeader.biClrUsed) - uintptr(l))

	// 分配缓冲区并读取DIB数据
	dataLen := uint32(dibHeader.biSizeImage)
	if dataLen == 0 {
		dataLen = dibHeader.biWidth * dibHeader.biHeight * uint32(dibHeader.biBitCount) / uint32(8)
	}
	dibData := make([]byte, dataLen)
	copy(dibData, (*[1 << 20]byte)(unsafe.Pointer(uintptr(l) + uintptr(offset)))[:dataLen])

	r, _, err := globalUnlock.Call(h)
	if r == 0 {
		_, _, _ = closeClipboard.Call()
		return []byte{}, err
	}

	closed, _, err := closeClipboard.Call()
	if closed == 0 {
		return []byte{}, err
	}

	// 创建图像
	img := image.NewRGBA(image.Rect(0, 0, int(dibHeader.biWidth), int(dibHeader.biHeight)))

	// 复制DIB数据到图像
	draw.Draw(
		img, img.Bounds(),
		&image.RGBA{
			Pix:    dibData,
			Stride: int(dibHeader.biWidth * uint32(dibHeader.biBitCount) / 8),
			Rect:   img.Bounds(),
		},
		image.Point{},
		draw.Src,
	)

	buf := bytes.NewBuffer(make([]byte, 0))
	png.Encode(buf, img)

	// 保存图片到文件
	file, err := os.Create("clipboard_image.png")
	if err != nil {
		log.Print("Failed to create file:", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Print("Failed to encode PNG:", err)
	}
	return []byte{}, nil
}
