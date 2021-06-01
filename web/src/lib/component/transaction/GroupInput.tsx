import React from 'react';
import {observer} from 'mobx-react-lite';
import {TextFieldProps} from '@material-ui/core';
import Autocomplete from '@material-ui/lab/Autocomplete';
import {RootStoreContext} from '../../store/RootStore';
import autocompleteProps from './autocompleteProps';
import IconInput from '../IconInput';
import DetailModel from 'lib/model/DetailModel';

interface IProps {
    detail: DetailModel;
}

const GroupInput = observer<IProps & Partial<TextFieldProps>, HTMLDivElement>(({detail, ...inputProps}, ref) => {
    const {groupStore} = React.useContext(RootStoreContext);
    return (
        <Autocomplete ref={ref} {...autocompleteProps} options={groupStore.groups} getOptionLabel={(g) => g.name}
            value={groupStore.getGroup(detail.transactionGroupId) ?? null}
            onChange={(_event, value) => detail.transactionGroupId = value ? value.id : undefined}
            filterOptions={(options, state) => options.filter((g) => g.name.toLowerCase().includes(state.inputValue.toLowerCase()))}
            renderInput={(params) => <IconInput {...params} {...inputProps} icon='workspaces' />} />
    );

}, {forwardRef: true});

export default GroupInput;
