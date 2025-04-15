import { ICategory, IReceipt } from '../types.ts';
import { faker } from '@faker-js/faker';

interface IInitialState {
    categories: ICategory[];
    receipts: IReceipt[];
    isLogin: boolean;
}
export const initialState: IInitialState = {
    isLogin: false,
    categories: [
        {
            id: '1',
            name: 'Еда 1',
            totalSum: 500,
        },
        {
            id: '2',
            name: 'Жилье 1',
            totalSum: 1200,
        },
        {
            id: '3',
            name: 'Развлечения 1',
            totalSum: 0,
        },

        {
            id: '4',
            name: 'ВБ 1',
            totalSum: 4333,
        },
        {
            id: '5',
            name: 'Командировка 1',
            totalSum: 542,
        },
        {
            id: '6',
            name: 'Метро 1',
            totalSum: 445,
        },

        {
            id: '7',
            name: 'Машина 1',
            totalSum: 52200,
        },
        {
            id: '8',
            name: 'Путешествия 1',
            totalSum: 4123,
        },
        {
            id: '9',
            name: 'Техника 1',
            totalSum: 47138,
        },
        {
            id: '10',
            name: 'Коммуналка 12',
            totalSum: 10000,
        },
        {
            id: '11',
            name: 'Еда3',
            totalSum: 500,
        },
        {
            id: '12',
            name: 'Жилье5',
            totalSum: 1200,
        },
        {
            id: '13',
            name: 'Развлечения4',
            totalSum: 0,
        },

        {
            id: '14',
            name: 'ВБ6',
            totalSum: 4333,
        },
        {
            id: '15',
            name: 'Командировка4',
            totalSum: 542,
        },
        {
            id: '16',
            name: 'Метро134',
            totalSum: 445,
        },

        {
            id: '17',
            name: 'Машина512',
            totalSum: 52200,
        },
        {
            id: '18',
            name: 'Путешествия54',
            totalSum: 4123,
        },
        {
            id: '19',
            name: 'Техника666',
            totalSum: 47138,
        },
        {
            id: '20',
            name: 'Коммуналка624',
            totalSum: 10000,
        },
    ],
    receipts: Array.from({ length: 20 }, (_, i) => ({
        id: faker.string.uuid(),
        name: `Чек ${i + 1} - ${faker.commerce.productName()}`,
        sum: parseFloat(faker.commerce.price({ min: 100, max: 10000 })),
        date: faker.date.between({
            from: new Date(2023, 0, 1),
            to: new Date(),
        }),
    })),
};
