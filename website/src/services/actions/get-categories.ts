import { apiService } from '../apiService.ts';
import { getErrorMessage } from '../../utils';
import { SET_CATEGORIES, SET_ERROR, SET_LOADING } from './index.ts';

export const getCategories = () => {
    const token = localStorage.getItem('authToken') || '';
    return async (dispatch: any) => {
        dispatch({
            type: SET_LOADING,
        });

        try {
            const response = await apiService.getReport(token);
            dispatch({
                type: SET_CATEGORIES,
                payload: Array.from(response.report).map(
                    ([name, totalSum]) => ({
                        id: name,
                        name,
                        totalSum,
                    })
                ),
            });
        } catch (error) {
            dispatch({
                type: SET_ERROR,
                payload: getErrorMessage(error),
            });
        }
    };
};
