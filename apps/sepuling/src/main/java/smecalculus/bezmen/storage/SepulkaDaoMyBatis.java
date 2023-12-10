package smecalculus.bezmen.storage;

import java.util.Optional;
import java.util.UUID;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import smecalculus.bezmen.core.StateDm.AggregateState;
import smecalculus.bezmen.core.StateDm.ExistenceState;
import smecalculus.bezmen.core.StateDm.PreviewState;
import smecalculus.bezmen.core.StateDm.TouchState;
import smecalculus.bezmen.storage.mybatis.SepulkaSqlMapper;

@RequiredArgsConstructor
public class SepulkaDaoMyBatis implements SepulkaDao {

    @NonNull
    private SepulkaStateMapper stateMapper;

    @NonNull
    private SepulkaSqlMapper sqlMapper;

    @Override
    public AggregateState add(@NonNull AggregateState state) {
        var stateEdge = stateMapper.toEdge(state);
        sqlMapper.insert(stateEdge);
        return state;
    }

    @Override
    public Optional<ExistenceState> getBy(@NonNull String externalId) {
        return sqlMapper.findByExternalId(externalId).map(stateMapper::toDomain);
    }

    @Override
    public Optional<PreviewState> getBy(@NonNull UUID internalId) {
        return sqlMapper.findByInternalId(internalId.toString()).map(stateMapper::toDomain);
    }

    @Override
    public void updateBy(TouchState state, UUID internalId) {
        var stateEdge = stateMapper.toEdge(state);
        var matchedCount = sqlMapper.updateBy(stateEdge, internalId.toString());
        if (matchedCount == 0) {
            throw new ContentionException();
        }
    }
}
