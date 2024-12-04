
# HertzDB

HertzDB is a lightweight search tool which currently supports token based search on indexed fields.

ToDo: Prefix and Term Search on Indexes


## Web API

Provides core REST APIs for essential features of HertzDB.

**Types**

**`SearchRequest`** : 

```
{
    "field":{
        "name":{name}
        "values":[{values}]
    }
}
```
**Endpoints**

**`POST :/v1/{index}/ingest`** : Take form-data with key file for ingesting file through multi-part, and index the same for the provided `index`.

**`GET :/v1/{index}/search`** : With body `SearchRequest`, return arrays of the matched documents.

## Core

Current functionalities include parsing and indexing json files for token based searching. This is done efficiently via batch processing the input file, 1000 documents at a time, and the indexed metadata is saved in a index file with extension `.hz` 

**`package index`**

**`index.DeserializeIndexDocumentMeta(filename string)`** : To be used to retrive index metadata

**`IndexDocument.ParseEntireFile(field[]string)`** : To be used to create token mapping for the provided json file


**`index.SearchTerm(filename string, field string, term string)`** : To be used to retrive location for documents for the given search term

**`index.GetDocument(filename string, locs [][2]int64)`** : To be used to retrive `JSON` documents from the provided file.
