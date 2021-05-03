import React from 'react';
import {observer} from 'mobx-react-lite';
import {TextFieldProps} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import {RootStoreContext} from '../../store/RootStore';
import {ITransactionDetail} from 'src/lib/model/TransactionModel';
import autocompleteProps from './autocompleteProps';
import {CategoryModel} from 'src/lib/model/CategoryModel';
import {AccountModel} from 'src/lib/model/account/AccountModel';
import IconInput from '../IconInput';

interface IProps {
    detail: ITransactionDetail;
}

export type IOption = CategoryModel | AccountModel;

const CategoryInput = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({detail, ...inputProps}, ref) => {
    const {accountStore, categoryStore} = React.useContext(RootStoreContext);
    const options = React.useMemo(() => {
        const categories = categoryStore.categories as IOption[];
        return categories.concat(accountStore.accounts);
    }, [accountStore.accounts, categoryStore.categories]);
    const selectedValue = detail.transactionCategoryId
        ? categoryStore.getCategory(detail.transactionCategoryId) ?? null
        : accountStore.getAccount(detail.relatedDetail?.transaction.accountId) ?? null;
    return (
        <Autocomplete ref={ref} {...autocompleteProps} options={options} getOptionLabel={(o) => o.displayName}
            value={selectedValue}
            filterOptions={(options, state) => options.filter((o) => o.displayName.toLowerCase().includes(state.inputValue.toLowerCase()))}
            renderInput={(params) => <IconInput {...params} {...inputProps} icon='category' />} />
    );
}, {forwardRef: true});

export default CategoryInput;