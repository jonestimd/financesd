import React from 'react';
import {observer} from 'mobx-react-lite';
import {TextField, TextFieldProps} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import {RootStoreContext} from '../../store/RootStore';
import {translate} from '../../i18n/localize';
import {ITransactionDetail} from 'src/lib/model/TransactionModel';
import autocompleteProps from './autocompleteProps';
import {CategoryModel} from 'src/lib/model/CategoryModel';
import {AccountModel} from 'src/lib/model/account/AccountModel';

interface IProps {
    detail: ITransactionDetail;
}

interface IOption {
    accountId?: number;
    categoryId?: number;
    label: string;
}

const categoryOption = (c: CategoryModel) => ({categoryId: c.id, label: c.displayName});
const accountOption = (a: AccountModel) => ({accountId: a.id, label: a.displayName});

const CategoryInput = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({detail, ...inputProps}, ref) => {
    const {accountStore, categoryStore} = React.useContext(RootStoreContext);
    const options = React.useMemo(() => {
        return categoryStore.categories.map<IOption>(categoryOption).concat(accountStore.accounts.map(accountOption));
    }, [accountStore.accounts, categoryStore.categories]);
    const selectedValue = detail.transactionCategoryId
        ? options.find((o) => o.categoryId === detail.transactionCategoryId) ?? null
        : options.find((o) => o.accountId === detail.relatedDetail?.transaction.accountId) ?? null;
    return (
        <Autocomplete ref={ref} {...autocompleteProps} options={options} getOptionLabel={(o) => o.label ?? '?'}
            value={selectedValue}
            filterOptions={(options, state) => options.filter((o) => o?.label.toLowerCase().includes(state.inputValue))}
            renderInput={(params) => <TextField {...params} {...inputProps} label={translate('Category')} />} />
    );
}, {forwardRef: true});

export default CategoryInput;
