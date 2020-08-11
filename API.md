# Example API using

[<= Back to main README file](README.md)

## 1. /api/v1/resize (for resizing image)
 
### Call parameters example
```json
{
  "user_id": "a493e097-6f4c-493d-9a82-e612b3d7e53d", 
  "width": 1400, 
  "height": 200, 
  "request_id":"zzz1"
}
+ multy part file (for example images.jpeg)
```
### Response example
```json
{
    "user_id": "a4Y93e097-6f4c-493d-9a82-e612b3d7e53d",
    "request_id": "zzz1",
    "err_code": 0,
    "err_msg": "",
    "image_id": 1941592313,
    "original_image_path": "https://amazonaws.com/a393e097-6f4c-493d-9a82-e612b3d7e53d/1941592313/images.jpeg",
    "resized_image_path": "https://amazonaws.com/a393e097-6f4c-493d-9a82-e612b3d7e53d/1941592313/images_1400x200.jpeg"
}
```
## 2. /api/v1/resize-by-id (for resizing image that previously was resized)

### Call parameters example
```json
{
	"image_id": 1941592313,
	"user_id":"a393e097-6f4c-493d-9a82-e612b3d7e53d",
	"request_id":"asdad",
	"width":22,
	"height":345
}
```
### Response example
```json
{
    "user_id": "a393e097-6f4c-493d-9a82-e612b3d7e53d",
    "request_id": "asdad",
    "err_code": 0,
    "err_msg": "",
    "image_id": 1941592313,
    "original_image_path": "https://amazonaws.com/a393e097-6f4c-493d-9a82-e612b3d7e53d/1941592313/images.jpeg",
    "resized_image_path": "https://amazonaws.com/a393e097-6f4c-493d-9a82-e612b3d7e53d/1941592313/images_22x345.jpeg"
}
```
## 3. /api/v1/list (show all processed images with path to original, resized images and resize params)

### Call parameters example
```text
user_id=a393e097-6f4c-493d-9a82-e612b3d7e53d&request_id=qq12
```
### Response example
```json
{
    "user_id": "a393e097-6f4c-493d-9a82-e612b3d7e53d",
    "request_id": "qq12",
    "err_code": 0,
    "err_msg": "",
    "data": [
        {
            "PicId": 3993602761,
            "Url": "https://amazonaws.com/a393e097-6f4c-493d-9a82-e612b3d7e53d/3993602761/ca760b70976b52578da88e06973af542.jpg",
            "ResizedImages": [
                {
                    "Url": "https://amazonaws.com/a393e097-6f4c-493d-9a82-e612b3d7e53d/3993602761/ca760b70976b52578da88e06973af542_1400x200.jpeg",
                    "Width": 1400,
                    "Height": 200
                },
                {
                    "Url": "https://amazonaws.com/a393e097-6f4c-493d-9a82-e612b3d7e53d/3993602761/ca760b70976b52578da88e06973af542_2212x345.jpeg",
                    "Width": 2212,
                    "Height": 345
                }
            ]
        },
        {
            "PicId": 1941592313,
            "Url": "https://amazonaws.com/a393e097-6f4c-493d-9a82-e612b3d7e53d/1941592313/images_resized_300_400.jpeg",
            "ResizedImages": [
                {
                    "Url": "https://amazonaws.com/a393e097-6f4c-493d-9a82-e612b3d7e53d/1941592313/images_resized_300_400_1400x200.jpeg",
                    "Width": 1400,
                    "Height": 200
                }
            ]
        }
    ]
}
```

## Error Codes
| Code| Description | 
| --- | --- |
| 600 | Empty request |
| 601 | File not found in request |
| 602 | Params not set in request |
| 603 | Cannot parse request params |
| 604 | Invalid values in request params |
| 605 | Cannot resize image |
| 606 | Cannot upload image to cloud store |
| 607 | Cannot save request results to DB |
| 608 | Cannot get user images from DB |
| 609 | Image not found |
| 610 | Cannot download file |
