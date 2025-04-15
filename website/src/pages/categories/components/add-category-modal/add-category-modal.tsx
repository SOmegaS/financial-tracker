import { TextInput } from '@mantine/core';
import { Modal } from '../../../../common-components';
import { useForm } from '@mantine/form';
import { memo } from 'react';

interface AddCategoryModalProps {
    isOpen: boolean;
    onClose: () => void;
}

const AddCategoryModal = ({ isOpen, onClose }: AddCategoryModalProps) => {
    const form = useForm({
        initialValues: {
            name: '',
        },
        validate: {
            name: (value) =>
                value.trim().length < 2 ? 'Название слишком короткое' : null,
        },
    });

    const handleSubmit = () => {
        console.log(form.values);
        form.reset();
        onClose();
    };

    if (!isOpen) return null;

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            onSave={handleSubmit}
            title="Добавить категорию"
            isDisabled={!form.isValid()}
            saveText="Добавить"
        >
            <TextInput
                label="Название категории"
                placeholder="Введите название"
                {...form.getInputProps('name')}
                mb="md"
            />
        </Modal>
    );
};

export default memo(AddCategoryModal);
