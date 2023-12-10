package smecalculus.bezmen.storage.mybatis

import org.apache.ibatis.annotations.Insert
import org.apache.ibatis.annotations.Param
import org.apache.ibatis.annotations.Select
import org.apache.ibatis.annotations.Update
import smecalculus.bezmen.storage.SepulkaStateEm

interface SepulkaSqlMapper {
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
        """,
    )
    fun insert(state: SepulkaStateEm.AggregateRoot)

    @Select(
        """
        SELECT
            internal_id as internalId,
            external_id as externalId
        FROM sepulkas
        WHERE external_id = #{externalId}
        """,
    )
    fun findByExternalId(externalId: String): SepulkaStateEm.Existence?

    @Select(
        """
        SELECT
            internal_id as internalId,
            external_id as externalId,
            created_at as createdAt
        FROM sepulkas
        WHERE internal_id = #{internalId}
        """,
    )
    fun findByInternalId(internalId: String): SepulkaStateEm.Preview?

    @Update(
        """
        UPDATE sepulkas
        SET revision = revision + 1,
            updated_at = #{state.updatedAt}
        WHERE internal_id = #{id}
        AND revision = #{state.revision}
        """,
    )
    fun updateBy(
        @Param("state") state: SepulkaStateEm.Touch,
        @Param("id") internalId: String,
    ): Int
}
