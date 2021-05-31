import React from 'react';
import {shallow} from 'enzyme';
import {RootStore} from '../store/RootStore';
import CompaniesDialog from './CompaniesDialog';
import {DialogActions, DialogTitle, IconButton} from '@material-ui/core';
import CompanyListModelClass from '../model/account/CompanyListModel';
import {newAccountModel, newCompanyModel} from 'test/accountFactory';
import AccountStore from '../store/AccountStore';
import CompanyRow from '../model/account/CompanyRow';
import {IColumn} from './table/Column';
import Table from './table/Table';
import TextCellEditor from './table/TextCellEditor';

jest.mock('../model/account/CompanyListModel');

describe('CompaniesDialog', () => {
    const rootStore = new RootStore();
    const onClose = jest.fn();
    const mockConstructor = CompanyListModelClass as jest.MockedClass<typeof CompanyListModelClass>;
    const CompanyListModel = jest.requireActual<{default: typeof CompanyListModelClass}>('../model/account/CompanyListModel').default;
    const model = new CompanyListModel({columns: 2}, {companies: [newCompanyModel()]} as AccountStore);

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
        mockConstructor.mockReturnValue(model);
    });
    describe('name column', () => {
        it('displays company name', () => {
            const component = shallow(<CompaniesDialog onClose={onClose} />);
            const column = component.find(Table).prop<IColumn<CompanyRow>[]>('columns')[0];

            expect(column.render(new CompanyRow({id: -1, name: 'the company'}))).toEqual('the company');
        });
        it('highlight changes', () => {
            const component = shallow(<CompaniesDialog onClose={onClose} />);
            const column = component.find(Table).prop<IColumn<CompanyRow>[]>('columns')[0];

            expect((column.className as (c: unknown) => string)({isChanged: true})).toEqual('changed');
            expect((column.className as (c: unknown) => string)({isChanged: false})).toEqual('');
        });
        it('is editable', () => {
            const component = shallow(<CompaniesDialog onClose={onClose} />);
            const column = component.find(Table).prop<IColumn<CompanyRow>[]>('columns')[0];

            const editor = column.editor!;

            const row = new CompanyRow({id: -1, name: 'the name'});
            expect(editor.Component).toEqual(TextCellEditor);
            expect(editor.getValue(row)).toEqual(row.name);
            editor.setValue(row, 'new name');
            expect(row.name).toEqual('new name');
        });
    });
    describe('accounts column', () => {
        it('displays account count', () => {
            const component = shallow(<CompaniesDialog onClose={onClose} />);
            const column = component.find(Table).prop<IColumn<CompanyRow>[]>('columns')[1];

            expect(column.render(new CompanyRow({id: -1, name: '', accounts: [newAccountModel()]}))).toEqual(1);
        });
    });
    describe('close button', () => {
        it('closes the dialog', () => {
            const component = shallow(<CompaniesDialog onClose={onClose} />);

            component.find(DialogTitle).find(IconButton).simulate('click');

            expect(onClose).toBeCalledTimes(1);
        });
    });
    describe('add button', () => {
        it('calls model.addCompany', () => {
            jest.spyOn(model, 'rows', 'get').mockReturnValue(2);
            jest.spyOn(model, 'addCompany');
            jest.spyOn(model, 'setEditCell');
            const component = shallow(<CompaniesDialog onClose={onClose} />);

            component.find(DialogActions).find(IconButton).at(0).simulate('click');

            expect(model.addCompany).toBeCalledTimes(1);
            expect(model.setEditCell).toBeCalledWith({row: 1, column: 0});
        });
    });
    describe('delete button', () => {
        it('is disabled for pending delete', () => {
            jest.spyOn(model, 'isDelete').mockReturnValue(true);

            const component = shallow(<CompaniesDialog onClose={onClose} />);

            expect(component.find(DialogActions).find(IconButton).at(1)).toBeDisabled();
        });
        it('is disabled for company with accounts', () => {
            jest.spyOn(model, 'selected', 'get').mockReturnValue(new CompanyRow({id: -1, name: '', accounts: [newAccountModel()]}));

            const component = shallow(<CompaniesDialog onClose={onClose} />);

            expect(component.find(DialogActions).find(IconButton).at(1)).toBeDisabled();
        });
        it('calls model.delete', () => {
            jest.spyOn(model, 'delete');
            const component = shallow(<CompaniesDialog onClose={onClose} />);

            component.find(DialogActions).find(IconButton).at(1).simulate('click');

            expect(model.delete).toBeCalledTimes(1);
        });
    });
    const saveButtonTests = [
        {name: 'is disabled if no changes', changed: false, valid: true, disabled: true},
        {name: 'is disabled if invalid', changed: true, valid: false, disabled: true},
        {name: 'is enabled if changed and valid', changed: true, valid: true, disabled: false}];
    describe('apply button', () => {
        saveButtonTests.forEach(({name, changed, valid, disabled}) => it(name, () => {
            jest.spyOn(model, 'isChanged', 'get').mockReturnValue(changed);
            jest.spyOn(model, 'isValid', 'get').mockReturnValue(valid);

            const component = shallow(<CompaniesDialog onClose={onClose} />);

            expect(component.find('#apply-companies')).toHaveProp('disabled', disabled);
        }));
        it('saves changes', () => {
            jest.spyOn(model, 'isChanged', 'get').mockReturnValue(true);
            jest.spyOn(model, 'isValid', 'get').mockReturnValue(true);
            jest.spyOn(model, 'save').mockResolvedValue(true);
            const component = shallow(<CompaniesDialog onClose={onClose} />);

            component.find('#apply-companies').simulate('click');

            expect(model.save).toBeCalledTimes(1);
        });
    });
    describe('save button', () => {
        saveButtonTests.forEach(({name, changed, valid, disabled}) => it(name, () => {
            jest.spyOn(model, 'isChanged', 'get').mockReturnValue(changed);
            jest.spyOn(model, 'isValid', 'get').mockReturnValue(valid);

            const component = shallow(<CompaniesDialog onClose={onClose} />);

            expect(component.find('#save-companies')).toHaveProp('disabled', disabled);
        }));
        it('saves changes and closes on success', async () => {
            jest.spyOn(model, 'isChanged', 'get').mockReturnValue(true);
            jest.spyOn(model, 'isValid', 'get').mockReturnValue(true);
            jest.spyOn(model, 'save').mockResolvedValue(true);
            const component = shallow(<CompaniesDialog onClose={onClose} />);

            await component.find('#save-companies').prop<() => Promise<void>>('onClick')();

            expect(model.save).toBeCalledTimes(1);
            expect(onClose).toBeCalledTimes(1);
        });
        it('saves changes and stays open on error', async () => {
            jest.spyOn(model, 'isChanged', 'get').mockReturnValue(true);
            jest.spyOn(model, 'isValid', 'get').mockReturnValue(true);
            jest.spyOn(model, 'save').mockResolvedValue(false);
            const component = shallow(<CompaniesDialog onClose={onClose} />);

            await component.find('#save-companies').prop<() => Promise<void>>('onClick')();

            expect(model.save).toBeCalledTimes(1);
            expect(onClose).not.toBeCalled();
        });
    });
});
