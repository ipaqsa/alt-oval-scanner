# Alt Oval Scanner (AOS)

## Overview
This is app for OVAL scanning, support only ALT Linux.

## Usage

The app provide three objects for scanning:

1. Host system
2. Image
3. Virtual Machine(Future)

Example: 

scan image:
```
aos image registry.altlinux.org/alt/alt -c example.yml
```

scan host:
```
aos host -c example.yml
```

## JSON Output Structure
The app can get output in json format (use `-o PATH_TO_JSON_FILE` flag)
```
{
    package: 
    [
        {
            title,
            version,
            installedVersion,
            references,
            CVEs,
            severity,
        },
    ]
} 

```