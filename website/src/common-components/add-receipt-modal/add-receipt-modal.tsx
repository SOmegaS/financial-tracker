import { TextInput, NumberInput } from '@mantine/core';
import {
    IconReceipt,
    IconCalendar,
    IconCurrencyRubel,
} from '@tabler/icons-react';
import { useForm } from '@mantine/form';
import { Modal } from '../../common-components';
import { DateInput } from '@mantine/dates';
import { memo } from 'react';

interface AddReceiptModalProps {
    opened: boolean;
    onClose: () => void;
    categoryName: string;
}

function AddReceiptModal({
    opened,
    onClose,
    categoryName,
}: AddReceiptModalProps) {
    const form = useForm({
        initialValues: {
            title: '',
            date: new Date(),
            amount: 0,
        },
        validate: {
            title: (value) =>
                value.trim().length < 2 ? 'Название слишком короткое' : null,
            amount: (value) =>
                value <= 0 ? 'Сумма должна быть больше 0' : null,
            date: (value) =>
                !value ? 'Дата обязательна для заполнения' : null,
        },
    });

    const handleSubmit = () => {
        console.log(form.values);
        form.reset();
        onClose();
    };

    return (
        <Modal
            title={`Добавить чек в ${categoryName}`}
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
