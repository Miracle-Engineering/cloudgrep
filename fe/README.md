## Prerequisites

Set AWS profile (credentials) in terminal. Export AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY

## Configuration settings and contstants

Majority of configuration constants are placed in "constants/globals.ts" file. Most of them are related to filtering and paging adjustmenst.
That file contains comments and description for constants (how and what for are they used)

## Install and run

Go to the FE project directory (fe) and run the following commands:

### `npm i` or `npm install`

Installs all needed dependency packages.

### `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

### `npm test`

Launches the test runner in the interactive watch mode.

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

## Continuous Integration

When creating a PR with changes in the `./fe` directory, a github action called [frontend-asset](/.github/workflows/frontend-asset.yml) is triggered to build the app for production and update the assets used by the Go app in `./static`. This action will create a commit to the existing PR to update the assets.
