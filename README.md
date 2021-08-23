# basset API Geo

API creted to manage geographical information. 

>This API is private for basset only.
>
>You can see a full public services documentation [here](https://docs.basset.ws/)

## Domain entities

### Countries

#### Useful endpoints
```bash
curl "https://internal.basset.ws/geo/countries?alpha2_code=AR"
curl "https://internal.basset.ws/geo/countries/59c2a51ad592ce268b9312cf"
```
#### Model
```json
{
    id: "ObjectId",
    name: {
        "lang":"string"
    },
    official_name: {
        "lang":"string"
    },
    alpha2_code: "string",
    alpha3_code: "string
}
```

### States

#### Useful endpoints

```bash
curl "https://internal.basset.ws/geo/states?country_id=59c2a51ad592ce268b9312cf"
curl "https://internal.basset.ws/geo/states/59cbc66bd592ce53955d2c09"
```

#### Model
```json
{
    id: "ObjectId",
    name: {
        "lang":"string"
    },
    country_id: "ObjectId",
    timezone: "string"
}
```

### Cities

#### Useful endpoints

```bash
curl "https://internal.basset.ws/geo/cities?country_id=59c2a51ad592ce268b9312cf"
curl "https://internal.basset.ws/geo/cities?state_id=59cbc66bd592ce53955d2c09"
curl "https://internal.basset.ws/geo/cities/59cbc66cd592ce53955d2c21"
```

#### Model
```json
{
    id: "ObjectId",
    name: {
        "lang":"string"
    },
    country_id: "ObjectId",
    state_id: "ObjectId",
}
```

### Neighbourhoods

#### Useful endpoints

```bash
curl "https://internal.basset.ws/geo/neighbourhoods?country_id=59c2a51ad592ce268b9312cf"
curl "https://internal.basset.ws/geo/neighbourhoods?state_id=59cbc66bd592ce53955d2c09"
curl "https://internal.basset.ws/geo/neighbourhoods?city_id=59cbc66cd592ce53955d2c21"
curl "https://internal.basset.ws/geo/neighbourhoods/59cbc66fd592ce53955d307f"
```

#### Model
```json
{
    id: "ObjectId",
    name: {
        "lang":"string"
    },
    country_id: "ObjectId",
    state_id: "ObjectId",
    city_id: "ObjectId",
}
```

### Airports

#### Useful endpoints

```bash
curl "https://internal.basset.ws/geo/airports?country_id=59c2a51ad592ce268b9312cf"
curl "https://internal.basset.ws/geo/airports?state_id=59cbc66bd592ce53955d2c09"
curl "https://internal.basset.ws/geo/airports?city_id=59cbc66cd592ce53955d2c21"
curl "https://internal.basset.ws/geo/airports/59dbd2eed592ce971a7cb2af"
```

#### Model
```json
{
    id: "ObjectId",
    name: {
        "lang":"string"
    },
    country_id: "ObjectId",
    state_id: "ObjectId",
    city_id: "ObjectId",
    iata_code:"string"
}
```
--- 

## Docker build

```bash
docker build --build-arg version=$version -t api-geo:$version .
```
## Docker Run local

```bash
docker run -d -p 80:8080 --rm -e env="-e $env" --name api-geo api-geo:$version
```

