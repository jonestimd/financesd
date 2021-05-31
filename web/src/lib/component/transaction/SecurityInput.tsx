import React from 'react';
import {observer} from 'mobx-react-lite';
import {TextFieldProps} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import {RootStoreContext} from '../../store/RootStore';
import TransactionModel from 'lib/model/TransactionModel';
import autocompleteProps from './autocompleteProps';
import IconInput from '../IconInput';

interface IProps {
    transaction: TransactionModel;
}

const SecurityInput = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({transaction, ...inputProps}, ref) => {
    const {securityStore} = React.useContext(RootStoreContext);
    return (
        <Autocomplete ref={ref} {...autocompleteProps} options={securityStore.securities} getOptionLabel={(s) => s.displayName}
            value={securityStore.getSecurity(transaction.securityId) ?? null}
            onChange={(_event, value) => transaction.securityId = value?.id ?? undefined}
            filterOptions={(options, state) => options.filter((s) => s.displayName.toLocaleLowerCase().includes(state.inputValue.toLowerCase()))}
            renderInput={(params) => <IconInput {...params} {...inputProps} icon='request_page' />} />
    );

}, {forwardRef: true});

export default SecurityInput;
