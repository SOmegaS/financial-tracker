import {
    LOGIN,
    LOGOUT,
    SET_CATEGORIES,
    SET_CREATE_RECEIPT_LOADING,
    SET_ERROR,
    SET_LOADING,
    SET_RECEIPTS,
    SET_TOKEN,
} from './services/actions';

export interface ICategory {
    id: string;
    name: string;
    totalSum: number;
}

export interface IReceipt {
    id: string;
    name: string;
    sum: number;
    date: string;
}

export interface IState {
    categories: ICategory[];
    receipts: IReceipt[];
    isAuth: boolean;
    isLoading: boolean;
    error: string | null;
    token: string | null;
    isCreateReceiptLoading: boolean;
}

export type Action =
    | { type: typeof SET_LOADING }
    | { type: typeof LOGIN; payload: string }
    | { type: typeof SET_TOKEN; payload: string }
    | { type: typeof LOGOUT }
    | { type: typeof SET_ERROR; payload: string }
    | { type: typeof SET_CREATE_RECEIPT_LOADING; payload: boolean }
    | {
          type: typeof SET_CATEGORIES;
          payload: ICategory[];
      }
    | { type: typeof SET_RECEIPTS; payload: IReceipt[] };
