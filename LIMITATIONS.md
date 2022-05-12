
# Know limitations and issues to fix

## Multi-regions support

The current implementations finds out the AWS region from the local configuration.  
But we should be able to fetch resources in different regions.  
Providers shouldn't be region-specific, but instead work across an entire cloud, but can be configured to filter on certain regions.  
Restricting providers to specific regions makes handling global services more complicated (such as AWS's IAM and route53) since we would need to work out which "region" we use to fetch those resources.

## Query on flatten properties name/value
Flatten the property values might create some limitations when filtering on resources. Here are a few examples below.
```js
//a user might want to search on SecurityGroups or SecurityGroups.GroupId
"name": "SecurityGroups[0][GroupId]",
"value": "sg-0c87bc57b6e994081"

"name": "SecurityGroups[1][GroupId]",
"value": "sg-0c87bc57b6e994081"
```

```js
//a user might want to search on IamInstanceProfile[Id] or IamInstanceProfile.Id but not Id
"name": "IamInstanceProfile[Id]",
"value": "AIPATTS7AQXJUXEPIKDD7"
```

```js
//a user might want to search on Placement[AvailabilityZone] or AvailabilityZone
"name": "Placement[AvailabilityZone]",
"value": "us-east-1b"
```

Options:
- We could store key alias names so these fields could be searched with alternate names.
- We could store expose and store the properties as a JSON object and use [PostgreSQL JSON Functions](https://www.postgresql.org/docs/12/functions-json.html)


## No pagination when fetching resources

Currently we query all the resource for one type in one call, we would need to revisit this to work with large AWS accounts.

Options:
- The AWS SDK offfers [Pagination Methods](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/making-requests.html)
- Have dedicated Go routines to fetch the data and to store it and use go Channel(s) to communicate between goroutines
- Implement the Iterator Pattern

## Returning resources without a tag is not supported

This is not supported at this time. Implementation TBD

Options:
- Add a query param "?exclude-tags=cluster,team": exlude tags `cluster` and `team`
- Generate a SQL query using EXCEPT such as:
    ```sql
    SELECT id FROM resources EXCEPT SELECT id FROM resources join tags where tag.name in ("cluster", "team")
    ```

