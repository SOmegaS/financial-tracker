import { TextInput, NumberInput, Autocomplete } from '@mantine/core';
import {
    IconReceipt,
    IconCalendar,
    IconCurrencyRubel,
    IconCategory,
} from '@tabler/icons-react';
import { useForm } from '@mantine/form';
import { Modal } from '../../common-components';
import { DateInput } from '@mantine/dates';
import { memo, useMemo, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { createReceipt } from '../../services/actions/create-receipt.ts';
import { IState } from '../../types.ts';
import { getCategories } from '../../services/actions/get-categories.ts';

interface AddReceiptModalProps {
    opened: boolean;
    onClose: () => void;
    categoryName?: string;
}

function AddReceiptModal({
                             opened,
                             onClose,
                             categoryName = '',
                         }: AddReceiptModalProps) {
    const [searchValue] = useState('');
    const categories = useSelector((state: IState) => state.categories) || [];
    const categoryOptions = useMemo(
        () => categories.map((category) => category.name),
        [categories]
    );

    const isLoading = useSelector(
        (state: IState) => state.isCreateReceiptLoading
    );
    const dispatch = useDispatch();

    const filteredCategories = useMemo(
        () =>
            categoryOptions.filter((name) =>
                name.toLowerCase().includes(searchValue.toLowerCase())
            ),
        [searchValue, categoryOptions]
    );

    const form = useForm({
        initialValues: {
            title: '',
            date: new Date(),
            amount: 0,
            category: categoryName || '',
        },
        validate: {
            title: (value) =>
                value.trim().length < 2 ? 'Название слишком короткое' : null,
            amount: (value) =>
                value <= 0 ? 'Сумма должна быть больше 0' : null,
            date: (value) =>
                !value ? 'Дата обязательна для заполнения' : null,
            category: (value) =>
                !categoryName && !value ? 'Выберите категорию' : null,
        },
    });

    const handleCreate = () => {
        onClose();
        // @ts-ignore

        dispatch(getCategories());
    };

    const handleSubmit = () => {
        // @ts-ignore
        dispatch(
            // @ts-ignore
            createReceipt(
                form.values.title,
                form.values.amount,
                form.values.category || categoryName,
                form.values.date,
                // @ts-ignore
                handleCreate
            )
        );
        form.reset();
    };

    if (!opened) return null;

    return (
        <Modal
            isLoading={isLoading}
            title={`Добавить чек${categoryName ? ` в ${categoryName}` : ''}`}
            isOpen={opened}
            saveText="Добавить"
            onClose={() => {
                form.reset();
                onClose();
            }}
            onSave={form.onSubmit(handleSubmit)}
            isDisabled={!form.isValid() || !form.isDirty()}
        >
            <form onSubmit={form.onSubmit(handleSubmit)}>
                {!categoryName && (
                    <Autocomplete
                        clearable
                        label="Категория"
                        placeholder="Выберите или введите новую категорию"
                        withAsterisk
                        mb="md"
                        leftSection={<IconCategory size={16} />}
                        data={filteredCategories}
                        value={form.values.category}
                        onChange={(value) =>
                            form.setFieldValue('category', value)
                        }
                    />
                )}
                <TextInput
                    label="Название чека"
                    placeholder="Например: Продукты"
                    withAsterisk
                    mb="md"
                    leftSection={<IconReceipt size={16} />}
                    {...form.getInputProps('title')}
                />
                <DateInput
                    label="Дата"
                    placeholder="Выберите дату"
                    withAsterisk
                    mb="md"
                    leftSection={<IconCalendar size={16} />}
                    valueFormat="DD.MM.YYYY"
                    clearable
                    maxDate={new Date()}
                    {...form.getInputProps('date')}
                />
                <NumberInput
                    label="Сумма"
                    placeholder="Введите сумму"
                    withAsterisk
                    mb="md"
                    leftSection={<IconCurrencyRubel size={16} />}
                    min={0}
                    decimalScale={2}
                    {...form.getInputProps('amount')}
                />
            </form>
        </Modal>
    );
}

export default memo(AddReceiptModal);
