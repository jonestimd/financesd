import React from 'react';
import {observer} from 'mobx-react-lite';
import {TextField, TextFieldProps} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import {RootStoreContext} from '../../store/RootStore';
import {translate} from '../../i18n/localize';
import TransactionModel from 'src/lib/model/TransactionModel';
import autocompleteProps from './autocompleteProps';

interface IProps {
    transaction: TransactionModel;
}

const SecurityInput = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({transaction, ...inputProps}, ref) => {
    const {securityStore} = React.useContext(RootStoreContext);
    return (
        <Autocomplete ref={ref} {...autocompleteProps} options={securityStore.securities} getOptionLabel={(s) => s.displayName}
            value={securityStore.getSecurity(transaction.securityId) ?? null}
            filterOptions={(options, state) => options.filter((s) => s.displayName.toLocaleLowerCase().includes(state.inputValue))}
            renderInput={(params) => <TextField {...params} {...inputProps} label={translate('Security')} />} />
    );

}, {forwardRef: true});

export default SecurityInput;
