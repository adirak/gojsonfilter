# gojsonfilter

This project using to filter, validate and initial value for map data using json configuration

**Example Data :**
```
{
    "firstName":"Supote",
    "lastName":"Sirimahayarn",
    "fullName": "Sirimahayarn-supote",
    "age": 42,
    "list": [
        {
            "a": "A",
            "b": "B",
            "c": "C"
        },
        {
            "a": "AA",
            "b": "BB",
            "c": "CC"
        }
    ]
}
```

**Example Json Filter :**
```
[
    {
        "default": "",
        "max": 0,
        "min": 0,
        "name": "fullName",
        "regexp": "",
        "required": true,
        "type": "string",
        "validated": true
    },
    {
        "all": false,
        "name": "list",
        "required": true,
        "type": "array",
        "validated": true,
        "min": 0,
        "max": 0,
        "children": [
            {
                "name": "[0]",
                "type": "map",
                "validated": true,
                "required": true,
                "min": 0,
                "max": 0,
                "children": [
                    {
                        "name": "a",
                        "type": "string",
                        "default": "",
                        "validated": true,
                        "required": false,
                        "min": 0,
                        "max": 0,
                        "regexp": ""
                    }
                ]
            }
        ]
    }
]

```

**Result Data :**
```
{
    "fullName": "Sirimahayarn-supote",
    "list": [
        {
            "a": "A"
        },
        {
            "a": "AA"
        }
    ]
}

```

**Configuration of filter:**

- name : string ==> field name to filter
- type : string ==> field type to filter
- required : bool ==> field cannot null
- validated : bool ==> turn on validate mode
- min : int ==> min value for integer type or min length for string type
- max : int ==> max value for integer type or max length for string type
- regexp : string ==> regular expression for validate string
- default : any ==> default value to initial if field is null
- children : array ==> child field of this field

