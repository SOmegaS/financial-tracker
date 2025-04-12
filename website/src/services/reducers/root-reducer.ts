import { initialState } from '../initial-state';

const rootReducer = (state = initialState, action: any) => {
    switch (action.type) {
        default: {
            return {
                ...state,
            };
        }
    }
};
export default rootReducer;
