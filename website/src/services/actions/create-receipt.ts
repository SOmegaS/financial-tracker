import { apiService } from '../apiService.ts';
import { getErrorMessage, getUserIdFromToken } from '../../utils';
import { SET_CREATE_RECEIPT_LOADING, SET_ERROR } from './index.ts';

export const createReceipt = (
    name: string,
    amount: number,
    category: string,
    date: Date,
    callback: () => {}
) => {
    const token = localStorage.getItem('authToken') || '';
    return async (dispatch: any) => {
        dispatch({
            type: SET_CREATE_RECEIPT_LOADING,
            payload: true,
        });

        try {
            await apiService.createBill(
                name,
                amount,
                category,
                getUserIdFromToken(token),
                Math.floor(date.getTime() / 1000),
                token
            );
        } catch (error) {
            dispatch({
                type: SET_ERROR,
                payload: getErrorMessage(error),
            });
        } finally {
            callback();
            dispatch({
                type: SET_CREATE_RECEIPT_LOADING,
                payload: false,
            });
        }
    };
};
