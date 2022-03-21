import React from 'react'
import ReactDOM from 'react-dom'
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import App from './App'
import {ThemeProvider} from "./Theme/Theme";


ReactDOM.render(
    <React.StrictMode>
        <ThemeProvider>
            <App/>
        </ThemeProvider>
    </React.StrictMode>,
    document.getElementById('root')
)
