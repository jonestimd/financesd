import {newSecurity, newSecurityModel} from 'test/securityFactory';
import {SecurityModel} from './SecurityModel';

describe('SecurityModel', () => {
    const security = newSecurity();

    describe('constructor', () => {
        it('populates security properties', () => {
            const model = new SecurityModel(security);

            expect(model).toEqual(security);
        });
    });
    describe('get displayName', () => {
        it('returns name', () => {
            expect(new SecurityModel(security).displayName).toEqual(security.name);
        });
        it('returns name and symbol', () => {
            const model = newSecurityModel({symbol: 'abc'});

            expect(model.displayName).toEqual(model.name + ' (abc)');
        });
    });
});
