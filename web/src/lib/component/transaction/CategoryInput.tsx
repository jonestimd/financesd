import React from 'react';
import {observer} from 'mobx-react-lite';
import {TextFieldProps} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import {RootStoreContext} from '../../store/RootStore';
import DetailModel from 'lib/model/DetailModel';
import autocompleteProps from './autocompleteProps';
import {CategoryModel} from 'lib/model/CategoryModel';
import {AccountModel} from 'lib/model/account/AccountModel';
import IconInput from '../IconInput';

interface IProps {
    detail: DetailModel;
}

export type IOption = CategoryModel | AccountModel;

const CategoryInput = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({detail, ...inputProps}, ref) => {
    const {accountStore, categoryStore} = React.useContext(RootStoreContext);
    const options = React.useMemo(() => {
        const categories = categoryStore.categories as IOption[];
        return categories.concat(accountStore.accounts);
    }, [accountStore.accounts, categoryStore.categories]);
    return (
        <Autocomplete ref={ref} {...autocompleteProps} options={options} getOptionLabel={(o) => o.displayName}
            value={detail.category} onChange={(_event, value) => detail.category = value}
            filterOptions={(options, state) => options.filter((o) => o.displayName.toLowerCase().includes(state.inputValue.toLowerCase()))}
            renderInput={(params) => <IconInput {...params} {...inputProps} icon='category' />} />
    );
}, {forwardRef: true});

export default CategoryInput;
