import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './normalize.css';
import '@mantine/core/styles.css';
import '@mantine/dates/styles.css';
import '@mantine/notifications/styles.css';
import App from './App.tsx';
import {
    createTheme,
    MantineProvider,
    TypographyStylesProvider,
} from '@mantine/core';
import { Provider } from 'react-redux';
import store from './services/store.ts';
import { Notifications } from '@mantine/notifications';

const theme = createTheme({
    fontFamily: 'Roboto, sans-serif',
    fontFamilyMonospace: 'Courier New, monospace',
    headings: { fontFamily: 'Roboto, sans-serif' },
});

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <Provider store={store}>
            <MantineProvider theme={theme}>
                <TypographyStylesProvider>
                    <Notifications />
                    <App />
                </TypographyStylesProvider>
            </MantineProvider>
        </Provider>
    </StrictMode>
);
