import { ICategory, IReceipt } from '../types.ts';
import { faker } from '@faker-js/faker';

interface IInitialState {
    categories: ICategory[];
    receipts: IReceipt[];
}
export const initialState: IInitialState = {
    categories: [
        {
            id: '1',
            name: 'Еда',
            totalSum: 500,
        },
        {
            id: '2',
            name: 'Жилье',
            totalSum: 1200,
        },
        {
            id: '3',
            name: 'Развлечения',
            totalSum: 0,
        },

        {
            id: '4',
            name: 'ВБ',
            totalSum: 4333,
        },
        {
            id: '5',
            name: 'Командировка',
            totalSum: 542,
        },
        {
            id: '6',
            name: 'Метро',
            totalSum: 445,
        },

        {
            id: '7',
            name: 'Машина',
            totalSum: 52200,
        },
        {
            id: '8',
            name: 'Путешествия',
            totalSum: 4123,
        },
        {
            id: '9',
            name: 'Техника',
            totalSum: 47138,
        },
        {
            id: '10',
            name: 'Коммуналка',
            totalSum: 10000,
        },
        {
            id: '11',
            name: 'Еда',
            totalSum: 500,
        },
        {
            id: '12',
            name: 'Жилье',
            totalSum: 1200,
        },
        {
            id: '13',
            name: 'Развлечения',
            totalSum: 0,
        },

        {
            id: '14',
            name: 'ВБ',
            totalSum: 4333,
        },
        {
            id: '15',
            name: 'Командировка',
            totalSum: 542,
        },
        {
            id: '16',
            name: 'Метро',
            totalSum: 445,
        },

        {
            id: '17',
            name: 'Машина',
            totalSum: 52200,
        },
        {
            id: '18',
            name: 'Путешествия',
            totalSum: 4123,
        },
        {
            id: '19',
            name: 'Техника',
            totalSum: 47138,
        },
        {
            id: '20',
            name: 'Коммуналка',
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
