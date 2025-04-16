import { LOGIN, SET_ERROR, SET_LOADING } from './index.ts';
import { apiService } from '../apiService.ts';
import { getErrorMessage } from '../../utils';
import { v4 as uuidv4 } from 'uuid';

export const register = (
    username: string,
    password: string,
    requestId = uuidv4()
) => {
    return async (dispatch: any) => {
        dispatch({
            type: SET_LOADING,
        });
        try {
            const response = await apiService.register(
                username,
                password,
                requestId
            );
            dispatch({
                type: LOGIN,
                payload: response.jwt,
            });
            localStorage.setItem('authToken', response.jwt);
        } catch (error) {
            dispatch({
                type: SET_ERROR,
                payload: getErrorMessage(error),
            });
        }
    };
};
