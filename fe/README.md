# CryptoBot Front-End Part

It's based on react-starter https://github.com/coryhouse/react-slingshot.git

## Available comands:
1. **npm start** - runs dev environment
1. **npm run build** - create build in dirrectory *./dist*

## Project structure
All code lays inside of **./src** folder.
**./components/** - folder with react components
**./styles/** - folder with SCSS
**./utils/config.js** - file with constants(e.g. ratesUpdateInterval)
**./utils/service.js** - place of integration with back-end server. There are api methods.
**./index.html** - entry point of application. Pay your attention to container 
`<div id="app"></div>`. *DON'T TOUCH THAT DIV.*
**./index.js** - root of react application.

## How to setup project
1. open project folder in terminal
1. complete **npm i** and wait until it will be finished
1. run any of available comands

## How to add new component
1. add this one to **./components/** folder.
2. attach it to any page
3. rebuild application(in dev mode it happens automaticaly)

## How to change styles
1. Just change styles in **./styles/** folder.
3. rebuild application(in dev mode it happens automaticaly)

## How to attach new library
1. Just use [npm](https://docs.npmjs.com/)
3. rebuild application(in dev mode it happens automaticaly, but you perhaps will need to restart dev mode)

## If you have any issues with build
FAQ: https://github.com/coryhouse/react-slingshot/blob/master/docs/FAQ.md

## Technologies
Slingshot offers a rich development experience using the following technologies:

| **Tech** | **Description** |**Learn More**|
|----------|-------|---|
|  [React](https://facebook.github.io/react/)  |   Fast, composable client-side components.    | [Pluralsight Course](https://www.pluralsight.com/courses/react-flux-building-applications) 
|  [React Router](https://github.com/reactjs/react-router) | A complete routing library for React | [Pluralsight Course](https://www.pluralsight.com/courses/react-flux-building-applications) |
|  [Babel](http://babeljs.io) |  Compiles ES6 to ES5. Enjoy the new version of JavaScript today.     | [ES6 REPL](https://babeljs.io/repl/), [ES6 vs ES5](http://es6-features.org), [ES6 Katas](http://es6katas.org), [Pluralsight course](https://www.pluralsight.com/courses/javascript-fundamentals-es6)    |
| [Webpack](https://webpack.js.org) | Bundles npm packages and our JS into a single file. Includes hot reloading via [react-transform-hmr](https://www.npmjs.com/package/react-transform-hmr). | [Quick Webpack How-to](https://github.com/petehunt/webpack-howto) [Pluralsight Course](https://www.pluralsight.com/courses/webpack-fundamentals)|
| [Browsersync](https://www.browsersync.io/) | Lightweight development HTTP server that supports synchronized testing and debugging on multiple devices. | [Intro vid](https://www.youtube.com/watch?time_continue=1&v=heNWfzc7ufQ) 
| [ESLint](http://eslint.org/)| Lint JS. Reports syntax and style issues. Using [eslint-plugin-react](https://github.com/yannickcr/eslint-plugin-react) for additional React specific linting rules. | |
| [SASS](http://sass-lang.com/) | Compiled CSS styles with variables, functions, and more. | [Pluralsight Course](https://www.pluralsight.com/courses/better-css)|
| [PostCSS](https://github.com/postcss/postcss) | Transform styles with JS plugins. Used to autoprefix CSS |
| [Editor Config](http://editorconfig.org) | Enforce consistent editor settings (spaces vs tabs, etc). | [IDE Plugins](http://editorconfig.org/#download) |
| [npm Scripts](https://docs.npmjs.com/misc/scripts)| Glues all this together in a handy automated build. | [Pluralsight course](https://www.pluralsight.com/courses/npm-build-tool-introduction), [Why not Gulp?](https://medium.com/@housecor/why-i-left-gulp-and-grunt-for-npm-scripts-3d6853dd22b8#.vtaziro8n)  |

