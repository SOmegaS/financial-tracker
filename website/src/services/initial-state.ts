import { IState } from '../types.ts';

export const initialState: IState = {
    isAuth: false,
    isLoading: false,
    error: null,
    token: null,
    categories: [],
    receipts: [],
    isCreateReceiptLoading: false,
};
