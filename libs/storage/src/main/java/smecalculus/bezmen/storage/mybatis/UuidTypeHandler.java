package smecalculus.bezmen.storage.mybatis;

import java.sql.CallableStatement;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Types;
import java.util.UUID;
import org.apache.ibatis.type.BaseTypeHandler;
import org.apache.ibatis.type.JdbcType;

/**
 * If the issue ever gets resolved <a href="https://github.com/spring-projects/spring-data-relational/issues/1648"/>
 * Then we can use String type with overridden JdbcType in the edge model {@link smecalculus.bezmen.storage.EdgeSide.AggregateState#internalId}.
 * So that type handler can be removed completely.
 */
public class UuidTypeHandler extends BaseTypeHandler<UUID> {

    @Override
    public void setNonNullParameter(PreparedStatement ps, int i, UUID parameter, JdbcType jdbcType)
            throws SQLException {
        ps.setObject(i, parameter.toString(), Types.OTHER);
    }

    @Override
    public UUID getNullableResult(ResultSet rs, String columnName) throws SQLException {
        var uuidString = rs.getString(columnName);
        return uuidString == null ? null : UUID.fromString(uuidString);
    }

    @Override
    public UUID getNullableResult(ResultSet rs, int columnIndex) throws SQLException {
        var uuidString = rs.getString(columnIndex);
        return uuidString == null ? null : UUID.fromString(uuidString);
    }

    @Override
    public UUID getNullableResult(CallableStatement cs, int columnIndex) throws SQLException {
        var uuidString = cs.getString(columnIndex);
        return uuidString == null ? null : UUID.fromString(uuidString);
    }
}
