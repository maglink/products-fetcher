# Products Fetcher

##### gRPC server which can fetch products list and then show their list

##

#### Methods:

### Fetch(url)

Fetch and store products list in MongoDB persistent storage.
 
Url must return CSV with a format:
`PRODUCT NAME;PRICE`


### List(limit, offset, order)

Return list with fields:
```
name, price, last_update, updates_count
```

##

### How to run
```
git clone github.com/maglink/products-fetcher
cd products-fetcher
make docker-run
```





