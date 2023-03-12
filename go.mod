module github.com/bbredesen/go-vk-samples

go 1.20

require (
	github.com/bbredesen/go-vk v0.0.0-20230304041510-33ed978b8bfd
	github.com/bbredesen/vkm v0.2.2
	github.com/bbredesen/win32-toolkit v0.0.0-20230303234304-25b01e7ba2d4
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20221017161538-93cebf72946b
	github.com/udhos/gwob v0.0.0-20200524213453-619810f75817
	golang.org/x/sys v0.6.0
)

require github.com/chewxy/math32 v1.10.1 // indirect

replace github.com/bbredesen/go-vk => ./vk
