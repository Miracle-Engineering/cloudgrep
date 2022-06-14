module.exports = async ({ context, core, require}) => {
    const fs = require('fs')

    let env
    if (context.eventName == "workflow_dispatch") {
        env = context.payload.inputs.environment
    } else if (context.eventName == "push") {
        env = "prod"
    } else if (context.eventName == "pull_request") {
        env = "dev"
    } else {
        throw "Unknown env for push event " + context.eventName
    }

    let accounts = JSON.parse(fs.readFileSync("integration/aws/accounts.json"))

    if (!accounts.hasOwnProperty(env)) {
        throw "Unknown account for env " + process.env.ENV
    }

    core.exportVariable('ENV', env)
    core.exportVariable('AWS_ACCOUNT_ID', accounts[env])
    core.exportVariable('TF_DIR', 'integration/aws/terraform/' + env + '-environment')
}
