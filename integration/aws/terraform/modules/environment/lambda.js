exports.handler = async function(event, context, callback) {
    console.log('Hello, logs!');
    return context.logStreamName
}
