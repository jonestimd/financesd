import {RootStore} from './RootStore';
import * as entityUtils from '../model/entityUtils';
import * as agent from '../agent';
import {loadingSecurities, query} from './SecurityStore';
import {newSecurity, newSecurityModel} from 'src/test/securityFactory';
import {SecurityModel} from '../model/SecurityModel';

describe('SecurityStore', () => {
    const {securityStore, messageStore} = new RootStore();

    beforeEach(() => {
        securityStore['securitiesById'].clear();
    });
    describe('get securities', () => {
        it('sorts by name', () => {
            const security = newSecurityModel();
            securityStore['securitiesById'].set(security.id, security);
            jest.spyOn(entityUtils, 'sortValuesByName');

            const securities = securityStore.securities;

            expect(securities).toEqual([security]);
            expect(entityUtils.sortValuesByName).toBeCalledWith(securityStore['securitiesById']);
        });
    });
    describe('getSecurty', () => {
        it('returns security for ID', () => {
            const security = newSecurityModel();
            securityStore['securitiesById'].set(security.id, security);

            expect(securityStore.getSecurity(security.id)).toBe(security);
            expect(securityStore.getSecurity(parseInt(security.id))).toBe(security);
        });
        it('returns undefined for unknown ID', () => {
            expect(securityStore.getSecurity('-99')).toBeUndefined();
        });
    });
    describe('loadSecurities', () => {
        beforeEach(() => {
            securityStore['loading'] = false;
            jest.spyOn(messageStore, 'addProgressMessage');
            jest.spyOn(messageStore, 'removeProgressMessage');
        });
        it('loads securities if securitiesById is empty', async () => {
            const security = newSecurity();
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {securities: [security]}});

            await securityStore.loadSecurities();

            expect(securityStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingSecurities);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingSecurities);
            expect(agent.graphql).toBeCalledWith('/finances/api/v1/graphql', query);
            expect(securityStore.securities).toStrictEqual([new SecurityModel(security)]);
        });
        it('does nothing is already loading', async () => {
            securityStore['loading'] = true;
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await securityStore.loadSecurities();

            expect(securityStore['loading']).toBe(true);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('does nothing is already loaded', async () => {
            const category = newSecurityModel();
            securityStore['securitiesById'].set(category.id, category);
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await securityStore.loadSecurities();

            expect(securityStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('logs error from graphql', async () => {
            const error = new Error('API error');
            jest.spyOn(agent, 'graphql').mockRejectedValue(error);
            jest.spyOn(console, 'error').mockImplementation(() => { });

            await securityStore.loadSecurities();

            expect(securityStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingSecurities);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingSecurities);
            expect(console.error).toBeCalledWith('error gettting securities', error);
        });
    });
});
