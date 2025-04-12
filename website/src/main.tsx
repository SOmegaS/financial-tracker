import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './normalize.css';
import '@mantine/core/styles.css';
import '@mantine/dates/styles.css';
import App from './App.tsx';
import {
    createTheme,
    MantineProvider,
    TypographyStylesProvider,
} from '@mantine/core';

const theme = createTheme({
    fontFamily: 'Roboto, sans-serif',
    fontFamilyMonospace: 'Courier New, monospace',
    headings: { fontFamily: 'Roboto, sans-serif' },
});

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <MantineProvider theme={theme}>
            <TypographyStylesProvider>
                <App />
            </TypographyStylesProvider>
        </MantineProvider>
    </StrictMode>
);
