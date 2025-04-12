import { configureStore } from '@reduxjs/toolkit';
import rootReducer from './reducers/root-reducer.ts';

export default configureStore({
    reducer: rootReducer,
});
