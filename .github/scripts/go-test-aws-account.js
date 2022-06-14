module.exports = async ({ context, core, exec, require }) => {
    const fs = require('fs')
    const accounts = JSON.parse(fs.readFileSync("integration/aws/accounts.json"))
    const tfDir = "integration/aws/terraform/"

    var env = "prod"
    if (context.eventName == "pull_request") {
        let changedFiles = []
        opts = {
            listeners: {
                stdline: (data) => {
                    changedFiles.push(data.toString())
                },
            },
        }

        let refspec = context.payload.pull_request.base.sha + '...HEAD'

        await exec.exec('git', ['diff', '--name-only', refspec], opts)

        let hasTerraform = changedFiles.some(file => file.startsWith(tfDir))
        if (hasTerraform) {
            env = "dev"
        }
    }

    if (!accounts.hasOwnProperty(env)) {
        throw "Unknown account for env '" + env + "'"
    }

    core.exportVariable('AWS_ACCOUNT_ID', accounts[env]);
}
