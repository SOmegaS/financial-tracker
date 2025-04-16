import { initialState } from '../initial-state';
import {
    LOGIN,
    LOGOUT,
    SET_CATEGORIES,
    SET_CREATE_RECEIPT_LOADING,
    SET_ERROR,
    SET_LOADING,
    SET_RECEIPTS,
    SET_TOKEN,
} from '../actions';
import { Action, IState } from '../../types.ts';

const rootReducer = (state: IState = initialState, action: Action) => {
    switch (action.type) {
        case SET_RECEIPTS:
            return {
                ...state,
                isLoading: false,
                receipts: action.payload,
            };
        case SET_CATEGORIES:
            return {
                ...state,
                isLoading: false,
                categories: action.payload,
            };
        case SET_LOADING:
            return {
                ...state,
                isLoading: true,
                error: null,
            };
        case LOGIN:
            return {
                isLoading: false,
                token: action.payload,
                isAuth: true,
                error: null,
            };
        case SET_TOKEN:
            return {
                ...state,
                token: action.payload,
            };
        case LOGOUT:
            return {
                ...state,
                token: null,
                isAuth: false,
            };
        case SET_ERROR:
            return {
                ...state,
                isLoading: false,
                error: action.payload,
            };
        case SET_CREATE_RECEIPT_LOADING:
            return {
                ...state,
                isCreateReceiptLoading: action.payload,
            };
        default: {
            return {
                ...state,
            };
        }
    }
};
export default rootReducer;
