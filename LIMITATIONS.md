
# Know limitations and issues to fix


## Query on flatten properties name/value
Flatten the property values might create some limitations when filtering on resources. Here are a few examples below.
```json
//a user might want to search on SecurityGroups or SecurityGroups.GroupId
"name": "SecurityGroups[0][GroupId]",
"value": "sg-0c87bc57b6e994081"

"name": "SecurityGroups[1][GroupId]",
"value": "sg-0c87bc57b6e994081"
```

```json
//a user might want to search on IamInstanceProfile[Id] or IamInstanceProfile.Id but not Id
"name": "IamInstanceProfile[Id]",
"value": "AIPATTS7AQXJUXEPIKDD7"
```

```json
//a user might want to search on Placement[AvailabilityZone] or AvailabilityZone
"name": "Placement[AvailabilityZone]",
"value": "us-east-1b"
```

Options:
- We could store key alias names so these fields could be searched with alternate names.
- We could store this information as JSON and use [PostgreSQL JSON Functions](https://www.postgresql.org/docs/12/functions-json.html)


## No pagination when fetching resources

Currently we query all the resource for one type in one call, we would need to revisit this to work with large AWS accounts.

Options:
- The AWS SDK offfers [Pagination Methods](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/making-requests.html)
- Have dedicated Go routines to fetch the data and to store it and use go Channel(s) to communicate between goroutines
- Implement the Iterator Pattern
