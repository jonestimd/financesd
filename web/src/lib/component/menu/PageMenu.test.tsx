import React, {TransitionEvent} from 'react';
import {shallow} from 'enzyme';
import PageMenu from './PageMenu';
import MainMenu from './MainMenu';
import CompaniesMenu from './CompaniesMenu';
import AccountsMenu from './AccountsMenu';
import {newCompanyModel} from 'test/accountFactory';

describe('PageMenu', () => {
    it('defaults to MainMenu', () => {
        const component = shallow(<PageMenu currentPage='menu.accounts' />);

        expect(component.find(MainMenu)).toExist();
        expect(component.find(CompaniesMenu)).not.toExist();
        expect(component.find(AccountsMenu)).not.toExist();
    });
    it('shows CompaniesMenu', () => {
        const component = shallow(<PageMenu currentPage='menu.accounts' />);

        component.find(MainMenu).prop('setDepth')(1);

        expect(component.find(MainMenu)).toExist();
        expect(component.find(CompaniesMenu)).toExist();
        expect(component.find(AccountsMenu)).not.toExist();
    });
    it('hides MainMenu after transition', () => {
        const component = shallow(<PageMenu currentPage='menu.accounts' />);
        component.find(MainMenu).prop('setDepth')(1);

        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);

        expect(component.find(MainMenu)).not.toExist();
        expect(component.find(CompaniesMenu)).toExist();
        expect(component.find(AccountsMenu)).not.toExist();
    });
    it('returns to MainMenu from CompaniesMenu', () => {
        const component = shallow(<PageMenu currentPage='menu.accounts' />);
        component.find(MainMenu).prop('setDepth')(1);
        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);

        component.find(CompaniesMenu).prop('onBack')();

        expect(component.find(MainMenu)).toExist();
        expect(component.find(CompaniesMenu)).toExist();
        expect(component.find(AccountsMenu)).not.toExist();
    });
    it('hides CompaniesMenu after transition to MainMenu', () => {
        const component = shallow(<PageMenu currentPage='menu.accounts' />);
        component.find(MainMenu).prop('setDepth')(1);
        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);
        component.find(CompaniesMenu).prop('onBack')();

        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);

        expect(component.find(MainMenu)).toExist();
        expect(component.find(CompaniesMenu)).not.toExist();
        expect(component.find(AccountsMenu)).not.toExist();
    });
    it('shows AccountsMenu', () => {
        const component = shallow(<PageMenu currentPage='menu.accounts' />);
        component.find(MainMenu).prop('setDepth')(1);
        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);

        component.find(CompaniesMenu).prop('selectCompany')(newCompanyModel());

        expect(component.find(MainMenu)).not.toExist();
        expect(component.find(CompaniesMenu)).toExist();
        expect(component.find(AccountsMenu)).toExist();
    });
    it('hides CompaniesMenu after transition to AccountsMenu', () => {
        const component = shallow(<PageMenu currentPage='menu.accounts' />);
        component.find(MainMenu).prop('setDepth')(1);
        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);
        component.find(CompaniesMenu).prop('selectCompany')(newCompanyModel());

        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);

        expect(component.find(MainMenu)).not.toExist();
        expect(component.find(CompaniesMenu)).not.toExist();
        expect(component.find(AccountsMenu)).toExist();
    });
    it('returns to CompaniesMenu from AccountsMenu', () => {
        const component = shallow(<PageMenu currentPage='menu.accounts' />);
        component.find(MainMenu).prop('setDepth')(1);
        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);
        component.find(CompaniesMenu).prop('selectCompany')(newCompanyModel());
        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);

        component.find(AccountsMenu).prop('onBack')();

        expect(component.find(MainMenu)).not.toExist();
        expect(component.find(CompaniesMenu)).toExist();
        expect(component.find(AccountsMenu)).toExist();
    });
    it('hides AccountsMenu after transition to CompaniesMenu', () => {
        const component = shallow(<PageMenu currentPage='menu.accounts' />);
        component.find(MainMenu).prop('setDepth')(1);
        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);
        component.find(CompaniesMenu).prop('selectCompany')(newCompanyModel());
        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);
        component.find(AccountsMenu).prop('onBack')();

        component.find('.menu').prop('onTransitionEnd')!({} as TransitionEvent);

        expect(component.find(MainMenu)).not.toExist();
        expect(component.find(CompaniesMenu)).toExist();
        expect(component.find(AccountsMenu)).not.toExist();
    });
});
