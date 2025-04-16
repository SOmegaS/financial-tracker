import { apiService } from '../apiService.ts';
import { getErrorMessage } from '../../utils';
import { SET_ERROR, SET_LOADING, SET_RECEIPTS } from './index.ts';
import { api } from '../../proto/generated/common.ts';
import Bill = api.Bill;

export const getBills = (category: string) => {
    const token = localStorage.getItem('authToken') || '';
    return async (dispatch: any) => {
        dispatch({
            type: SET_LOADING,
        });

        try {
            const response = await apiService.getBills(token, category);
            dispatch({
                type: SET_RECEIPTS,
                payload: response.bills.map((bill: Bill) => ({
                    id: bill.name,
                    name: bill.name,
                    sum: bill.amount,
                    date: new Date(bill.ts.seconds * 1000).toLocaleDateString(
                        'ru-RU',
                        {
                            day: '2-digit',
                            month: '2-digit',
                            year: 'numeric',
                        }
                    ),
                })),
            });
        } catch (error) {
            dispatch({
                type: SET_ERROR,
                payload: getErrorMessage(error),
            });
        }
    };
};
