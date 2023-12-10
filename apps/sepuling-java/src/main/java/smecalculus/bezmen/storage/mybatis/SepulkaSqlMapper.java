package smecalculus.bezmen.storage.mybatis;

import java.util.Optional;
import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;
import org.apache.ibatis.annotations.Update;
import smecalculus.bezmen.storage.SepulkaStateEm.AggregateRoot;
import smecalculus.bezmen.storage.SepulkaStateEm.Existence;
import smecalculus.bezmen.storage.SepulkaStateEm.Preview;
import smecalculus.bezmen.storage.SepulkaStateEm.Touch;

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
    void insert(AggregateRoot state);

    @Select(
            """
            SELECT
                internal_id as internalId
            FROM sepulkas
            WHERE external_id = #{externalId}
            """)
    Optional<Existence> findByExternalId(String externalId);

    @Select(
            """
            SELECT
                external_id as externalId,
                created_at as createdAt
            FROM sepulkas
            WHERE internal_id = #{internalId}
            """)
    Optional<Preview> findByInternalId(String internalId);

    @Update(
            """
            UPDATE sepulkas
            SET revision = revision + 1,
                updated_at = #{state.updatedAt}
            WHERE internal_id = #{id}
            AND revision = #{state.revision}
            """)
    int updateBy(@Param("state") Touch state, @Param("id") String internalId);
}
