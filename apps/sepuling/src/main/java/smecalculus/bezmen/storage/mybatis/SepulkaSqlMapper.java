package smecalculus.bezmen.storage.mybatis;

import java.util.Optional;
import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;
import org.apache.ibatis.annotations.Update;
import smecalculus.bezmen.storage.StateEm.AggregateState;
import smecalculus.bezmen.storage.StateEm.ExistenceState;
import smecalculus.bezmen.storage.StateEm.PreviewState;
import smecalculus.bezmen.storage.StateEm.TouchState;

public interface SepulkaSqlMapper {

    @Insert(
            """
            INSERT INTO sepulkas (
                internal_id,
                external_id,
                revision,
                created_at,
                updated_at
            )
            VALUES (
                #{internalId},
                #{externalId},
                #{revision},
                #{createdAt},
                #{updatedAt}
            )
            """)
    void insert(AggregateState state);

    @Select(
            """
            SELECT
                internal_id as internalId
            FROM sepulkas
            WHERE external_id = #{externalId}
            """)
    Optional<ExistenceState> findByExternalId(String externalId);

    @Select(
            """
            SELECT
                external_id as externalId,
                created_at as createdAt
            FROM sepulkas
            WHERE internal_id = #{internalId}
            """)
    Optional<PreviewState> findByInternalId(String internalId);

    @Update(
            """
            UPDATE sepulkas
            SET revision = revision + 1,
                updated_at = #{state.updatedAt}
            WHERE internal_id = #{id}
            AND revision = #{state.revision}
           """)
    int updateBy(@Param("state") TouchState state, @Param("id") String internalId);
}
