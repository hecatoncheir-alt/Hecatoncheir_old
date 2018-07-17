# Hecatoncheir

<br>

> #### TODO:
> [✓] Language property <br>
> [0] Log <br>

---
Send message to broker:

```json
 {
 	"Message": "Need products of category of company",
 	"Data": {
               "Language":"en",
               "Company":{  
                  "ID":"0x2786",
                  "Name":"M.Video",
                  "IRI":"http://www.mvideo.ru/"
               },
               "Category":{  
                  "ID":"",
                  "Name":"Test category of M.Video company"
               },
               "City":{  
                  "ID":"0x2788",
                  "Name":"Москва"
               },
               "PageInstruction":{  
                  "uid":"0x2789",
                  "path":"smartfony-i-svyaz/smartfony-205",
                  "pageInPaginationSelector":".pagination-list .pagination-item",
                  "previewImageOfSelector":".product-tile-picture-link img",
                  "pageParamPath":"/f/page=",
                  "cityParamPath":"?cityId=",
                  "itemSelector":".grid-view .product-tile",
                  "nameOfItemSelector":".product-tile-title",
                  "linkOfItemSelector":".product-tile-title a",
                  "cityInCookieKey":"",
                  "cityIdForCookie":"",
                  "priceOfItemSelector":".product-price-current"
              }
    }
 }
```

Response for all connected clients:
```json
{  
   "Data":{  
      "Name":"Смартфон Samsung Galaxy S8 64Gb Черный бриллиант",
      "IRI":"http://www.mvideo.ru//products/smartfon-samsung-galaxy-s8-64gb-chernyi-brilliant-30027818",
      "PreviewImageLink":"img.mvideo.ru/Pdb/30027818m.jpg",
      "Language":"en",
      "Price":{  
         "Value":46990,
         "City":{  
            "ID":"0x2788",
            "Name":"Москва"
         },
         "DateTime":"2018-02-10T08:34:35.6055814Z"
      },
      "Company":{  
         "ID":"0x2786",
         "Name":"M.Video",
         "IRI":"http://www.mvideo.ru/"
      },
      "Category":{  
         "ID":"",
         "Name":"Test category of M.Video company"
      },
      "City":{  
         "ID":"0x2788",
         "Name":"Москва"
      }
   },
   "Message":"Product of category of company ready"
}
```

## Setup
Need NSQ:
```docker
docker pull nsqio/nsq

docker run --name lookupd -p 4160:4160 -p 4161:4161 nsqio/nsq /nsqlookupd

docker run --name nsqd -p 4150:4150 -p 4151:4151 nsqio/nsq /nsqd --broadcast-address=192.168.99.100 --lookupd-tcp-address=192.168.99.100:4160

#admin panel
docker run -p 4171:4171 nsqio/nsq /nsqadmin --lookupd-http-address=192.168.99.100:4161
```

Set NSQD addres:
```go
broker := broker.New()
	go broker.Connect("192.168.99.100", 4150)
```
