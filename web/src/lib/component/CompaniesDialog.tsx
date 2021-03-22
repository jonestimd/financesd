import React from 'react';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../store/RootStore';
import Table from './table/Table';
import {IColumn} from './table/Column';
import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Icon, IconButton, Typography} from '@material-ui/core';
import {translate} from '../i18n/localize';
import TextCellEditor from './table/TextCellEditor';
import CompanyRow from '../model/account/CompanyRow';
import CompanyListModel, {nameKey} from '../model/account/CompanyListModel';

const columns: IColumn<CompanyRow>[] = [{
    key: nameKey,
    className: (c) => c?.isChanged ? 'changed' : '',
    render: (c) => c.name,
    editor: {
        Component: TextCellEditor,
        getValue: (c) => c.name,
        setValue: (c, name) => c.setName(name),
    },
},{
    key: 'company.accounts',
    render: (c) => c.accounts,
    className: 'number',
}];

interface IProps {
    onClose: () => void;
}

const CompaniesDialog: React.FC<IProps> = ({onClose}) => {
    const {accountStore} = React.useContext(RootStoreContext);
    const model = React.useMemo(() => new CompanyListModel({headerSelector: 'thead', columns: columns.length}, accountStore), [accountStore]);
    const addCompany = () => {
        model.addCompany();
        model.setEditCell({row: model.rows - 1, column: 0});
    };
    const saveAndClose = async () => {
        if (await model.save()) onClose();
    };
    const disableSave = !model.isChanged || !model.isValid;
    const disableDelete = !model.selected || model.isDelete() || model.selected.accounts > 0;
    return (
        <Dialog open fullWidth id='company-list' className='modal' PaperProps={{square: true}} onClose={onClose}>
            <DialogTitle disableTypography>
                <Typography>{translate('companies.title')}</Typography>
                <IconButton size='small' onClick={onClose}><Icon>close</Icon></IconButton>
            </DialogTitle>
            <DialogContent>
                <Table columns={columns} data={model.items} selection={model} />
            </DialogContent>
            <DialogActions>
                <IconButton size='small' onClick={addCompany}><Icon>add</Icon></IconButton>
                <IconButton size='small' disabled={disableDelete} onClick={() => model.delete()}>
                    <Icon>delete</Icon>
                </IconButton>
                <Button id='apply-companies' variant='outlined' disabled={disableSave} onClick={() => model.save()}>Apply</Button>
                <Button id='save-companies' color='primary' variant='contained' disabled={disableSave} onClick={saveAndClose}>Save</Button>
            </DialogActions>
        </Dialog>
    );
};

export default observer(CompaniesDialog);
